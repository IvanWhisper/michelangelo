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
	timeoutCtx, cancel := context.WithTimeout(a.GetContext(), 2*time.Minute)
	defer cancel()
	app := exec.CommandContext(timeoutCtx, a.cmd, args...)
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	if err := app.Start(); err != nil {
		return err
	}
	defer app.Process.Kill()
	mlog.InfoCtx(fmt.Sprintf("PID[%d]%s Exec %v", app.Process.Pid, a.name, args), a.GetContext())
	completedCh := make(chan CompleteResult, 0)
	go func() {
		defer close(completedCh)
		if err := app.Wait(); err != nil {
			mlog.ErrorCtx(fmt.Sprintf("PID[%d]%s Exec %v Error %s", app.Process.Pid, a.name, args, err), a.GetContext())
			completedCh <- CompleteResult{Success: false, Error: err}
			return
		}
		completedCh <- CompleteResult{Success: true, Error: nil}
	}()
	select {
	case <-timeoutCtx.Done():
		return errors.New(fmt.Sprintf("PID[%d]%s Exec %v timeOut %d", app.Process.Pid, a.name, args, a.timeOut))
	case c := <-completedCh:
		if c.Success {
			return nil
		} else {
			return c.Error
		}
	}
}
