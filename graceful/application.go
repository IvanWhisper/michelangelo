package graceful

import (
	"context"
	"errors"
	"fmt"
	mlog "github.com/IvanWhisper/michelangelo/log"
	"os"
	"os/exec"
	"time"
)

type IApplication interface {
	SetName(value string)
	GetName() string

	SetCmd(value string)
	GetCmd() string

	SetWorkPath(value string)
	GetWorkPath() string

	SetTimeOut(value time.Duration)
	GetTimeOut() time.Duration

	SetContext(ctx context.Context)
	GetContext() context.Context

	Run(args ...string) error
}

type application struct {
	name     string
	cmd      string
	workPath string
	timeOut  time.Duration
	ctx      context.Context
}

type CompleteResult struct {
	Success bool
	Error   error
}

func (a *application) SetName(value string) {
	a.name = value
}

func (a *application) GetName() string {
	return a.name
}

func (a *application) SetCmd(value string) {
	a.cmd = value
}

func (a *application) GetCmd() string {
	return a.cmd
}

func (a *application) SetWorkPath(value string) {
	a.workPath = value
}

func (a *application) GetWorkPath() string {
	return a.workPath
}

func (a *application) SetTimeOut(value time.Duration) {
	a.timeOut = value
}

func (a *application) GetTimeOut() time.Duration {
	return a.timeOut
}

func (a *application) SetContext(ctx context.Context) {
	a.ctx = ctx
}

func (a *application) GetContext() context.Context {
	return a.ctx
}

func (a *application) Run(args ...string) error {
	timeoutCtx, cancel := context.WithTimeout(a.GetContext(), a.GetTimeOut())
	defer cancel()
	app := exec.CommandContext(timeoutCtx, a.cmd, args...)
	app.Dir = a.GetWorkPath()
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	if err := app.Start(); err != nil {
		return err
	}
	defer app.Process.Kill()
	completedCh := make(chan CompleteResult, 0)
	go func() {
		defer close(completedCh)
		if err := app.Wait(); err != nil {
			mlog.ErrorCtx(a.GetContext(), fmt.Sprintf("PID[%d]%s Exec %s %v %s", app.Process.Pid, a.name, a.cmd, args, err))
			completedCh <- CompleteResult{Success: false, Error: err}
			return
		}
		completedCh <- CompleteResult{Success: true, Error: nil}
	}()
	select {
	case <-timeoutCtx.Done():
		return errors.New(fmt.Sprintf("PID[%d]%s Exec %v timeOut %fs", app.Process.Pid, a.GetName(), args, a.GetTimeOut().Seconds()))
	case c := <-completedCh:
		if c.Success {
			return nil
		} else {
			return c.Error
		}
	}
}
