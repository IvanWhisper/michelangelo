package generate

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"

	mlog "github.com/IvanWhisper/michelangelo/infrastructure/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type TestMock struct {
	MockDB    *xorm.Engine
	MockCache *redis.Client
}

var env *TestMock

func GetMock() *TestMock {
	if env == nil {
		mlog.New(nil)
		e, err := xorm.NewEngine("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8mb4")
		if err != nil {
			panic(err)
		}
		e.ShowSQL(true)
		e.SetMaxIdleConns(10)
		e.SetMaxOpenConns(100)
		err = e.Sync2(new(IDRecord))
		if err != nil {
			panic(err)
		}
		c := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		env = &TestMock{MockDB: e, MockCache: c}
	}
	return env
}

func TestGenerater_TakeSingle(t *testing.T) {
	m := GetMock()
	g := &Generator{
		enableCache:  true,
		maxCacheStep: 10000,
		db:           m.MockDB,
		cache:        m.MockCache,
	}
	ctx := context.TODO()
	key := "test"
	takeCount := int64(10)
	takeStep := int64(10)
	start, _ := g.Cursor(ctx, key)
	var wg sync.WaitGroup
	for i := takeCount; i > 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Take(ctx, key, takeStep)
		}()
	}
	wg.Wait()
	end, _ := g.Cursor(ctx, key)
	target := start + takeCount*takeStep
	log.Println(fmt.Sprintf("end is %d，target is %d", end, target))
	if end != target {
		panic("no pass")
	}
}

func TestGenerater_TakeMore(t *testing.T) {
	m := GetMock()
	g := &Generator{
		enableCache:  true,
		maxCacheStep: 100,
		db:           m.MockDB,
		cache:        m.MockCache,
	}
	ctx := context.TODO()

	key1 := "mtest1"
	takeCount := int64(800)
	takeStep := int64(3)

	start, _ := g.Take(ctx, key1, 1)
	var wg sync.WaitGroup
	for i := takeCount; i > 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, _ := g.Take(ctx, key1, takeStep)
			log.Println("take id is " + strconv.FormatInt(id, 10))
		}()
	}
	wg.Wait()
	end, _ := g.Cursor(ctx, key1)
	target := start + takeCount*takeStep
	log.Println(fmt.Sprintf("end is %d，target is %d", end, target))
	if end != target {
		t.Error("no pass")
	}
}
