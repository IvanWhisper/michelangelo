package compense

import (
	"encoding/json"
	"errors"
	"github.com/IvanWhisper/michelangelo/dependencies"
	"time"

	"github.com/google/uuid"
)

var _ = dependencies.IOutTasker(&Task{})

type Task struct {
	ID         string `xorm:"varchar(64) pk 'Id'"`       // 主键
	TimeStamp  int64  `xorm:"bigint 'TimeStamp'"`        // 时间戳
	Group      int32  `xorm:"int index 'Group'"`         // 时间戳
	Name       string `xorm:"varchar(32) index 'Code'"`  // 时间戳
	Args       string `xorm:"text 'Args'"`               // 时间戳
	Retries    int64  `xorm:"int 'Retries'"`             // 时间戳
	Locker     string `xorm:"varchar(64) 'Locker'"`      // 时间戳
	LastError  string `xorm:"varchar(255) 'LastError'"`  // 时间戳
	LockExpire int64  `xorm:"bigint index 'LockExpire'"` // 时间戳
	ExecTime   int64  `xorm:"bigint 'ExecTime'"`         // 时间戳
}

func (p *Task) TableName() string {
	return "compensation_task"
}

func (p *Task) GetID() string {
	return p.ID
}

func (p *Task) GetName() string {
	return p.Name
}

func (p *Task) GetArgs() string {
	return p.Args
}

func (p *Task) GetParams() (map[string]interface{}, error) {
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

func (p *Task) GetTimeout() time.Duration {
	now := time.Now().Unix()
	if p.LockExpire > now {
		return time.Second * time.Duration(p.LockExpire-now)
	}
	return 0
}

func (p *Task) GetDeadline() int64 {
	return p.LockExpire
}

func (p *Task) GetOwner() string {
	return p.Locker
}

func (p *Task) GetLastError() string {
	return p.LastError
}

func (p *Task) GetRetries() int64 {
	return p.Retries
}

func (p *Task) GetGroup() int32 {
	return p.Group
}

func NewTask(timeout, margin time.Duration, group int32, name, args string) *Task {
	now := time.Now()
	expire := now.Add(margin).Add(timeout) // add extra 5 second
	return &Task{
		ID:         uuid.New().String(),
		TimeStamp:  now.Unix(),
		Group:      group,
		Name:       name,
		Args:       args,
		Retries:    0,
		LockExpire: expire.Unix(),
		Locker:     uuid.NewString(),
	}
}
