package compense

import (
	"context"
	"fmt"
	"github.com/IvanWhisper/michelangelo/dependencies"
	"sync"
	"time"

	"github.com/google/uuid"

	mlog "github.com/IvanWhisper/michelangelo/infrastructure/log"

	"github.com/go-redis/redis/v8"
	"xorm.io/xorm"
)

var _ = dependencies.ICompensator(&Compensation{})

type Compensation struct {
	db           *xorm.Engine
	cache        *redis.Client
	margin       time.Duration // lockExpire margin
	funcMap      map[string]dependencies.Exec
	execCallBack func(string)
	newAsyncCtx  func(ctx context.Context) context.Context
}

func (r *Compensation) SetDB(db *xorm.Engine) error {
	r.db = db
	err := db.Sync2(new(Task))
	if err != nil {
		return err
	}
	return nil
}

func (r *Compensation) SetNewAsyncCtx(v func(ctx context.Context) context.Context) error {
	r.newAsyncCtx = v
	return nil
}

func (r *Compensation) SetExecCallBack(v func(string)) error {
	r.execCallBack = v
	return nil
}

func (r *Compensation) Register(name string, exec dependencies.Exec) error {
	if r.funcMap == nil {
		r.funcMap = make(map[string]dependencies.Exec)
	}
	if _, ok := r.funcMap[name]; ok {
		return fmt.Errorf("%s func has registed", name)
	}
	r.funcMap[name] = exec
	return nil
}

func (r *Compensation) Execute(ctx context.Context, task dependencies.ITasker, timeout time.Duration) (int64, dependencies.IOutTasker, error) {
	t := NewTask(timeout, r.margin, task.GetGroup(), task.GetName(), task.GetArgs())
	rows, err := r.add(ctx, t)
	if err != nil {
		return rows, t, err
	}
	err = r.execute(ctx, t, timeout)
	return rows, t, err
}
func (r *Compensation) ExecuteAsync(ctx context.Context, task dependencies.ITasker, timeout time.Duration) (int64, dependencies.IOutTasker, error) {
	t := NewTask(timeout, r.margin, task.GetGroup(), task.GetName(), task.GetArgs())
	rows, err := r.add(ctx, t)
	if err != nil {
		return rows, t, err
	}
	asyncCtx := r.newAsyncCtx(ctx)
	go func() {
		err = r.execute(asyncCtx, t, timeout)
	}()
	return rows, t, err
}

func (r *Compensation) TaskByID(ctx context.Context, id string) (dependencies.IOutTasker, error) {
	session := r.db.NewSession()
	defer session.Close()
	session.Context(ctx)
	t := &Task{}
	b, err := session.ID(id).Get(t)
	if !b {
		return nil, err
	}
	return t, err
}

func (r *Compensation) Executes(ctx context.Context, group, limit int32, timeout time.Duration) error {
	tasks, err := r.Tasks(ctx, group, limit, timeout)
	if err != nil {
		return err
	}
	if len(tasks) < 1 {
		return nil
	}
	var wg sync.WaitGroup
	for _, v := range tasks {
		wg.Add(1)
		task := v
		asyncCtx := r.newAsyncCtx(ctx)
		go func() {
			defer wg.Done()
			err := r.execute(asyncCtx, task, timeout)
			mlog.InfoCtx(ctx, fmt.Sprintf("%s %s %s", task.GetID(), task.GetName(), err))
		}()
	}
	wg.Wait()
	return nil
}

func (r *Compensation) execute(ctx context.Context, task dependencies.ITasker, timeout time.Duration) (err error) {
	defer func() {
		rows, afterErr := r.after(ctx, task, err)
		r.notify(fmt.Sprintf("After Exec Rows %d Err%s", rows, afterErr))
	}()
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	ch := make(chan struct{})
	errCh := make(chan error)
	defer cancel()
	go func(subCtx context.Context) {
		defer func() {
			if errRecover := recover(); errRecover != nil {
				errCh <- fmt.Errorf("panic %s", errRecover)
			}
			close(ch)
			close(errCh)
		}()
		if v, ok := r.funcMap[task.GetName()]; ok {
			err := v(subCtx, task)
			if err != nil {
				errCh <- err
			} else {
				ch <- struct{}{}
			}
			r.notify(fmt.Sprintf("%s[%s] Err %s ", task.GetName(), task.GetArgs(), err))
		}
	}(timeoutCtx)

	select {
	case <-ch:
		return nil
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return fmt.Errorf("timeout %dms", timeout.Milliseconds())
	}
}

func (r *Compensation) Tasks(ctx context.Context, group, limit int32, timeout time.Duration) ([]dependencies.IOutTasker, error) {
	tasks, err := r.queryByType(ctx, group, limit, time.Now().Add(r.margin).Add(timeout).Unix()) // add extra 5 second
	if err != nil {
		return nil, err
	}
	result := make([]dependencies.IOutTasker, 0)
	for _, task := range tasks {
		v := task
		result = append(result, &v) //nolint:gosec
	}
	return result, nil
}

func (r *Compensation) Exec(ctx context.Context, exec func(*xorm.Session) (int64, error)) (int64, error) {
	session := r.db.NewSession()
	defer session.Close()
	session.Context(ctx)
	res, err := exec(session)
	mlog.InfoCtx(ctx, fmt.Sprintf("Exec res:%d,err:%s", res, err))
	return res, err
}

func (r *Compensation) after(ctx context.Context, task dependencies.ITasker, err error) (int64, error) {
	if err != nil {
		return r.recordLastError(ctx, task.GetID(), task.GetOwner(), err.Error())
	}
	return r.delete(ctx, task.GetID(), task.GetOwner())
}

func (r *Compensation) notify(message string) {
	if r.execCallBack != nil {
		r.execCallBack(message)
	}
}

func (r *Compensation) add(ctx context.Context, tasks ...*Task) (int64, error) {
	beans := make([]interface{}, 0)
	for _, task := range tasks {
		beans = append(beans, task)
	}
	res, err := r.Exec(ctx, func(session *xorm.Session) (int64, error) {
		return session.Insert(beans...)
	})
	mlog.InfoCtx(ctx, fmt.Sprintf("Exec res:%d,err:%s", res, err))
	return res, err
}

func (r *Compensation) queryByType(ctx context.Context, group, limit int32, lockExpire int64) (tasks []Task, err error) {
	session := r.db.NewSession()
	session.Context(ctx)
	defer func() {
		if p := recover(); p != nil {
			mlog.ErrorCtx(ctx, fmt.Sprintf("ExecWithTrans recover:%s", p))
			// mlog.Error(helper.GetErrorStack())
			session.Rollback()
		} else if err != nil {
			session.Rollback()
		}
		session.Close()
	}()
	if err := session.Begin(); err != nil {
		return nil, err
	}
	locker := uuid.NewString()
	rows, err := session.
		Cols("Locker", "Retries", "LockExpire").
		Where("`Group` = ? AND `LockExpire` <? AND `Locker`=''", group, time.Now().Unix()).
		Asc("`TimeStamp`").
		Limit(int(limit), 0).
		SetExpr("`Retries`", "Retries+1").
		Update(&Task{LockExpire: lockExpire, Locker: locker})
	if err != nil {
		return nil, err
	}
	tasks = make([]Task, 0)
	err = session.Where("Locker=?", locker).Find(&tasks)
	if err != nil {
		return tasks, err
	}
	mlog.InfoCtx(ctx, fmt.Sprintf("ExecWithTrans res:%d,err:%s", rows, err))
	if err := session.Commit(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Compensation) recordLastError(ctx context.Context, id, locker, msg string) (int64, error) {
	return r.Exec(ctx, func(session *xorm.Session) (int64, error) {
		return session.ID(id).Where("Locker=?", locker).Cols("ExecTime", "LastError", "Locker").Update(&Task{LastError: msg, ExecTime: time.Now().Unix()})
	})
}

func (r *Compensation) delete(ctx context.Context, id, locker string) (int64, error) {
	return r.Exec(ctx, func(session *xorm.Session) (int64, error) {
		return session.ID(id).Where("Locker=?", locker).Delete(&Task{})
	})
}
