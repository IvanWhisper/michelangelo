package generate

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	mlog "github.com/IvanWhisper/michelangelo/infrastructure/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type Generator struct {
	enableCache  bool
	maxCacheStep int64
	db           *xorm.Engine
	cache        *redis.Client
	keyPrefix    string
}

func NewGenerator(db *xorm.Engine, cache *redis.Client, enableCache bool, maxCacheStep int64, prefix string) func() (interface{}, error) {
	return func() (interface{}, error) {
		g := &Generator{
			enableCache:  enableCache,
			maxCacheStep: 10000,
			db:           db,
			cache:        cache,
			keyPrefix:    prefix,
		}
		if maxCacheStep > 0 {
			g.maxCacheStep = maxCacheStep
		}
		return g, nil
	}
}

type TakeFuncHandler func(ctx context.Context) (int64, error)

func FuncWithTimeout(h TakeFuncHandler, duration time.Duration) TakeFuncHandler {
	return func(ctx context.Context) (int64, error) {
		newctx, cancel := context.WithTimeout(ctx, duration)
		defer cancel()

		processDone := make(chan int64)
		defer close(processDone)
		errorDone := make(chan error)
		defer close(errorDone)

		go func() {
			r, e := h(newctx)
			if e != nil {
				errorDone <- e
			} else {
				processDone <- r
			}
		}()

		select {
		case <-ctx.Done():
			return 0, errors.New("action is timeout")
		case r := <-processDone:
			return r, nil
		case e := <-errorDone:
			return 0, e
		}
	}
}

func (g *Generator) TakeWithTimeout(ctx context.Context, key string, step int64, duration time.Duration) (int64, error) {
	return FuncWithTimeout(func(ctx context.Context) (int64, error) {
		return g.Take(ctx, key, step)
	}, duration)(ctx)
}

func (g *Generator) Take(ctx context.Context, key string, step int64) (int64, error) {
	mlog.DebugCtx(ctx, "Step01 Take Begin")
	if !g.enableCache || step >= g.maxCacheStep {
		mlog.DebugCtx(ctx, fmt.Sprintf("Step01.01 unable Cache[%v] or step[%d]>maxCacheStep[%d]", g.enableCache, step, g.maxCacheStep))
		if idr, err := g.dbIncrease(key, step); err != nil {
			mlog.Error(fmt.Sprintf("Generator:key[%s] %s", key, err))
			return 0, err
		} else {
			return idr.CurValue, nil
		}
	} else {
		mlog.DebugCtx(ctx, "Step01.01 enable Cache")
		if result, err := g.cacheIncrease(ctx, key, step, g.maxCacheStep); err != nil {
			return 0, err
		} else {
			if result == nil {
				str := fmt.Sprintf("Generator:key[%s]cacheIncrease %s", key, err)
				mlog.DebugCtx(ctx, "Step02 cacheIncrease result is nil "+str)
				return 0, errors.New(str)
			}
			if result.CurValue == -1 {
				mlog.DebugCtx(ctx, "Step02 cacheIncrease CurValue is -1")
				if idr, err := g.increase(ctx, key, step); err != nil {
					return 0, err
				} else {
					if idr == nil {
						return 0, fmt.Errorf("Generator:key[%s]increase result is err", key)
					}
					return idr.CurValue, err
				}
			} else {
				return result.CurValue, nil
			}
		}
	}

}

func (g *Generator) Cursor(ctx context.Context, key string) (int64, error) {
	if res, err := g.cacheCursor(ctx, key); err != nil {
		return 0, err
	} else {
		return res.CurValue, nil
	}
}
func (g *Generator) Cursors(ctx context.Context) (map[string]int64, error) {
	return nil, nil
}

// 临界时从数据库中提取批量冲入缓存，然后从缓存中取号
func (g *Generator) increase(ctx context.Context, key string, step int64) (*IDRecord, error) {
	mlog.DebugCtx(ctx, "Step01 increase begin")
	var fn = func() (interface{}, error) {
		// 在锁内再次尝试取号
		if result, err := g.cacheIncrease(ctx, key, step, g.maxCacheStep); err != nil {
			return result, err
		} else {
			if result == nil {
				return nil, fmt.Errorf("Generator:key[%s]cacheIncrease result is err", key)
			}
			if result.CurValue == -1 {
				// 从数据库提取号段
				dbs, err := g.dbIncrease(key, g.maxCacheStep)
				if err != nil {
					return nil, err
				}
				if dbs == nil {
					return nil, fmt.Errorf("Generator:key[%s]dbIncrease result is err", key)
				}
				mlog.DebugCtx(ctx, "Step02 increase dbIncrease "+dbs.Name+" value is "+strconv.FormatInt(dbs.CurValue, 10))
				// 写入缓存
				err = g.cacheBuildSource(ctx, key, dbs.CurValue-g.maxCacheStep)
				if err != nil {
					return nil, err
				}
				// 读取
				return g.cacheIncrease(ctx, key, step, g.maxCacheStep)
			} else {
				return result, nil
			}
		}
	}

	if r, err := LockInvoke(fn, key, time.Second*5); err != nil {
		return nil, err
	} else {
		if r == nil {
			return nil, nil
		}
		return r.(*IDRecord), nil
	}
}

func (g *Generator) cacheIncrease(ctx context.Context, name string, step int64, maxstep int64) (*IDRecord, error) {
	res := &IDRecord{}
	result := g.cache.Eval(ctx, LuaStInc, []string{g.buildKey(name, "source"), g.buildKey(name, "num")}, []string{strconv.Itoa(int(step)), strconv.Itoa(int(maxstep))})
	if err := result.Err(); err != nil {
		return nil, err
	}
	v := result.Val().(string)
	r := strings.Split(v, ",")
	res.Name = name
	value, _ := strconv.Atoi(r[1])
	res.CurValue = int64(value)
	return res, nil
}

func (g *Generator) cacheBuildSource(ctx context.Context, name string, start int64) error {
	mlog.DebugCtx(ctx, "cacheBuildSource begin")
	//result := g.cache.Set(ctx, g.buildKey(name, "source"), start, -1)
	result, err := g.cache.Eval(ctx, LuaBuildSource, []string{g.buildKey(name, "source")}, strconv.FormatInt(start, 10)).Result()
	if err != nil {
		return err
	} else {
		mlog.DebugCtx(ctx, fmt.Sprintf("cacheBuildSource set result is %s", result))
	}
	// 如果存在了key就不去修改值，不存在是设置初始值
	b := g.cache.SetNX(ctx, g.buildKey(name, "num"), start, -1)
	if b == nil {
		return errors.New("cache setnx result is nil")
	}
	return nil
}

func (g *Generator) cacheCursor(ctx context.Context, name string) (*IDRecord, error) {
	res := g.cache.Get(ctx, g.buildKey(name, "num"))
	if res != nil {
		if id, err := res.Int64(); err != nil {
			return nil, err
		} else {
			return &IDRecord{Name: name, CurValue: id}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("cache:key[%s]result is nil", name))
}

func (g *Generator) cacheCursors() ([]IDRecord, error) {
	return nil, nil
}

func (g *Generator) dbIncrease(name string, step int64) (*IDRecord, error) {
	p := &IDRecord{}
	session := g.db.NewSession()
	defer session.Close()
	// add Begin() before any action
	if err := session.Begin(); err != nil {
		// if returned then will rollback automatically
		return nil, err
	}
	if _, err := session.Exec(SQLInOrUp, name, step, step); err != nil {
		return nil, err
	}
	if _, err := session.Where("`Name`=?", name).Get(p); err != nil {
		return nil, err
	}
	if err := session.Commit(); err != nil {
		return nil, err
	}
	return p, nil
}

func (g *Generator) dbCursor(name string) (*IDRecord, error) {
	p := &IDRecord{}
	if _, err := g.db.Where("`Name`=?", name).Get(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (g *Generator) dbCursors() ([]IDRecord, error) {
	p := make([]IDRecord, 0)
	if err := g.db.Find(p); err != nil {
		return nil, err
	}
	res := make([]IDRecord, 0)
	for _, v := range p {
		res = append(res, v)
	}
	return res, nil
}

func (g *Generator) buildKey(key, suffix string) string {
	return g.keyPrefix + key + ":" + suffix
}

func LockInvoke(fn func() (interface{}, error), key string, timeout time.Duration) (interface{}, error) {
	v, b := loker_map.LoadOrStore(key, new(sync.Mutex))
	if !b {
		mlog.Debug("LockInvoke:locker is new, key is " + key)
	}
	v.(*sync.Mutex).Lock()
	println(key, " locker is lock")
	defer func() {
		// TODO:这里释放容易并发出问题,loker_map的大小一直递增，后期可以考虑使用一个锁池来控制，目前无限增长
		v.(*sync.Mutex).Unlock()
		mlog.Debug("LockInvoke:" + key + " locker is release")
	}()

	// 超时保障
	ch := make(chan interface{})
	defer close(ch)
	go func() {
		res, _ := fn()
		ch <- res
	}()
	select {
	case res := <-ch:
		return res, nil
	case <-time.After(timeout):
		fmt.Println("locker is timeout!!!")
		return nil, nil
	}
}
