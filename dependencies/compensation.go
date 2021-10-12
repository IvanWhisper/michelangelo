package dependencies

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"xorm.io/xorm"
)

var _ = ITasker(&DefaultTask{})

type ICompensator interface {
	SetDB(db *xorm.Engine) error
	SetExecCallBack(v func(string)) error
	SetNewAsyncCtx(v func(ctx context.Context) context.Context) error
	Register(name string, exec Exec) error
	Execute(ctx context.Context, task ITasker, timeout time.Duration) (int64, IOutTasker, error)
	ExecuteAsync(ctx context.Context, task ITasker, timeout time.Duration) (int64, IOutTasker, error)
	Executes(ctx context.Context, group, limit int32, timeout time.Duration) error
	TaskByID(ctx context.Context, id string) (IOutTasker, error)
}

type IExecutor interface {
	Registry(name string, exec Exec)
}

type Exec func(ctx context.Context, tasker ITasker) error

type ITasker interface {
	GetID() string
	GetGroup() int32
	GetName() string
	GetArgs() string
	GetParams() (map[string]interface{}, error)
	GetDeadline() int64
	GetOwner() string
}

type IOutTasker interface {
	ITasker
	GetLastError() string
	GetRetries() int64
}

type DefaultTask struct {
	ID       string
	Group    int32
	Name     string
	Args     string
	Deadline int64
	Owner    string
	LastErr  string
}

func (p *DefaultTask) GetID() string {
	return p.ID
}

func (p *DefaultTask) GetName() string {
	return p.Name
}

func (p *DefaultTask) GetParams() (map[string]interface{}, error) {
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(p.GetArgs()), &params)
	if err != nil {
		return params, err
	}
	if len(params) < 1 {
		return params, errors.New("params not enough")
	}
	return params, nil
}

func (p *DefaultTask) GetGroup() int32 {
	return p.Group
}

func (p *DefaultTask) GetArgs() string {
	return p.Args
}

func (p *DefaultTask) GetDeadline() int64 {
	return p.Deadline
}

func (p *DefaultTask) GetOwner() string {
	return p.Owner
}
