package compense

import (
	"context"
	"fmt"
	"github.com/IvanWhisper/michelangelo/dependencies"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"

	mlog "github.com/IvanWhisper/michelangelo/infrastructure/log"
)

func TestCompensation_Register(t *testing.T) {
	mlog.New(nil)
	c := MockCompensation()
	err := c.Register("a", func(ctx context.Context, tasker dependencies.ITasker) error {
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	err = c.Register("a", func(ctx context.Context, tasker dependencies.ITasker) error {
		return nil
	})
	if err == nil {
		t.Error("register same name no error")
	}
}

func TestCompensation_New_Success(t *testing.T) {
	mlog.New(nil)
	c := MockCompensation()
	ctx := context.Background()
	tA := &dependencies.DefaultTask{
		Name:  "A",
		Args:  "随便吧",
		Owner: "A",
	}
	_, result, _ := c.Execute(ctx, tA, time.Second*5)
	res, err := c.TaskByID(ctx, result.GetID())
	if err != nil {
		t.Error(err)
	}
	if res != nil {
		t.Error("res != nil")
	}
}

func TestCompensation_New_Error(t *testing.T) {
	mlog.New(nil)
	c := MockCompensation()
	ctx := context.Background()
	tB := &dependencies.DefaultTask{
		Name:  "B",
		Args:  "随便吧",
		Owner: "aA",
	}
	_, tb, _ := c.Execute(ctx, tB, time.Second*5)
	res, err := c.TaskByID(ctx, tb.GetID())
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("res is nil")
	}
	compareTask(tB, res, t.Error)
}

func TestCompensation_New_Panic(t *testing.T) {
	mlog.New(nil)
	c := MockCompensation()
	ctx := context.Background()
	tC := &dependencies.DefaultTask{
		Name:  "C",
		Args:  "随便吧",
		Owner: "XC",
	}
	_, tc, _ := c.Execute(ctx, tC, time.Second*5)
	res, err := c.TaskByID(ctx, tc.GetID())
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("res is nil")
	}
	compareTask(tC, res, t.Error)
}

func TestCompensation_New_Timeout(t *testing.T) {
	mlog.New(nil)
	c := MockCompensation()
	ctx := context.Background()
	tD := &dependencies.DefaultTask{
		Name:  "D",
		Args:  "随便吧",
		Owner: "XD",
	}
	_, td, _ := c.Execute(ctx, tD, time.Second*5)
	res, err := c.TaskByID(ctx, td.GetID())
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("res is nil")
	}
	compareTask(tD, res, t.Error)
}

func TestCompensation_Executes(t *testing.T) {
	mlog.New(nil)
	c := MockCompensation()
	ctx := context.Background()
	t1 := &dependencies.DefaultTask{
		Name:  "C",
		Args:  "随便吧",
		Owner: "XD",
	}
	t2 := &dependencies.DefaultTask{
		Name:  "B",
		Args:  "B随便吧",
		Owner: "B",
	}
	_, out1, _ := c.Execute(ctx, t1, time.Second*1)
	_, out2, _ := c.Execute(ctx, t2, time.Second*1)
	time.Sleep(time.Second * 7)
	c.Executes(ctx, 0, 1000, time.Second*1)
	res, err := c.TaskByID(ctx, out1.GetID())
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("res is nil")
	}
	if res.GetRetries() != 1 {
		t.Error("retries != 1")
	}
	compareTask(t1, res, t.Error)

	res, err = c.TaskByID(ctx, out2.GetID())
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("res is nil")
	}
	if res.GetRetries() != 1 {
		t.Error("retries != 1")
	}
	compareTask(t2, res, t.Error)
}

func compareTask(source dependencies.ITasker, result dependencies.ITasker, callback func(args ...interface{})) {
	if source.GetName() != result.GetName() {
		callback(fmt.Sprintf("[%s] name %s %s", result.GetID(), source.GetName(), result.GetName()))
	}
	if source.GetArgs() != result.GetArgs() {
		callback(fmt.Sprintf("[%s]args %s %s", result.GetID(), source.GetArgs(), result.GetArgs()))
	}
}

func MockCompensation() *Compensation {
	listener := make(chan string)
	go func() {
		defer close(listener)
		for v := range listener {
			println(v)
		}
	}()
	funcMap := make(map[string]dependencies.Exec)
	funcMap["A"] = func(ctx context.Context, tasker dependencies.ITasker) error {
		_, err := (&StructA{}).FuncA(ctx, tasker.GetArgs())
		return err
	}
	funcMap["B"] = func(ctx context.Context, tasker dependencies.ITasker) error {
		_, err := FuncB(ctx, tasker.GetArgs())
		return err
	}
	funcMap["C"] = func(ctx context.Context, tasker dependencies.ITasker) error {
		_, err := FuncC(ctx, tasker.GetArgs())
		return err
	}
	funcMap["D"] = func(ctx context.Context, tasker dependencies.ITasker) error {
		_, err := FuncD(ctx, tasker.GetArgs())
		return err
	}
	mock := GetMock()
	c := &Compensation{
		db:      mock.MockDB,
		cache:   mock.MockCache,
		funcMap: funcMap,
		margin:  time.Second * 5,
		execCallBack: func(s string) {
			println(s)
		},
		newAsyncCtx: func(ctx context.Context) context.Context {
			return context.Background()
		},
	}
	return c
}

type StructA struct {
	Name string
}

func (s *StructA) FuncA(ctx context.Context, args string) (string, error) {
	time.Sleep(time.Millisecond * 20)
	return "A", nil
}

func FuncB(ctx context.Context, args string) (string, error) {
	time.Sleep(time.Millisecond * 50)
	return "B", fmt.Errorf("%s exec failed", "B")
}

func FuncC(ctx context.Context, args string) (string, error) {
	time.Sleep(time.Millisecond * 10)
	panic("TMD出错了")
}

func FuncD(ctx context.Context, args string) (string, error) {
	time.Sleep(time.Second * 10)
	fmt.Println("10s 之后")
	panic(" 我不可能出现")
}

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
		err = e.Sync2(new(Task))
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
