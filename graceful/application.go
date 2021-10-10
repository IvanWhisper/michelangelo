package graceful

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	mlog "github.com/IvanWhisper/michelangelo/log"
	"os"
	"os/exec"
	"time"
)

var _ = IApplication(&application{})

type IApplication interface {
	SetName(value string)
	GetName() string

	SetCmd(value string)
	GetCmd() string

	SetPrintCh(value chan string)
	GetPrintCh() chan string

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
	printCh  chan string
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

func (a *application) SetPrintCh(value chan string) {
	a.printCh = value
}

func (a *application) GetPrintCh() chan string {
	return a.printCh
}

func (a *application) SetWorkPath(value string) {
	a.workPath = value
}

func (a *application) GetWorkPath() string {
	if a.workPath == "" {
		return "./"
	}
	return a.workPath
}

func (a *application) SetTimeOut(value time.Duration) {
	a.timeOut = value
}

func (a *application) GetTimeOut() time.Duration {
	if a.timeOut.Seconds() == 0 {
		return 10 * time.Second
	}
	return a.timeOut
}

func (a *application) SetContext(ctx context.Context) {
	a.ctx = ctx
}

func (a *application) GetContext() context.Context {
	if a.ctx != nil {
		return context.Background()
	}
	return a.ctx
}

func (a *application) initPrinter(app *exec.Cmd) error {
	if a.GetPrintCh() == nil {
		app.Stdout = os.Stdout
		app.Stderr = os.Stderr
	} else {
		output, err := app.StdoutPipe()
		app.Stderr = app.Stdout
		if err != nil {
			return err
		}
		go func() {
			defer close(a.GetPrintCh())
			for {
				tmp := make([]byte, 2048)
				_, err := output.Read(tmp)
				index := bytes.IndexByte(tmp, 0)
				a.GetPrintCh() <- string(tmp[0:index])
				if err != nil {
					break
				}
			}
		}()
	}
	return nil
}

func (a *application) Run(args ...string) error {
	timeoutCtx, cancel := context.WithTimeout(a.GetContext(), a.GetTimeOut())
	defer cancel()
	app := exec.CommandContext(timeoutCtx, a.cmd, args...) //nolint:gosec
	app.Dir = a.GetWorkPath()
	err := a.initPrinter(app)
	if err != nil {
		return err
	}
	if err := app.Start(); err != nil {
		return err
	}
	defer func() {
		if !app.ProcessState.Exited() {
			err := app.Process.Kill()
			if err != nil {
				mlog.ErrorCtx(a.GetContext(), err.Error())
			}
		}
	}()
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
		errMsg := fmt.Sprintf("PID[%d]%s Exec %v timeOut %fs", app.Process.Pid, a.GetName(), args, a.GetTimeOut().Seconds())
		mlog.InfoCtx(a.GetContext(), errMsg)
		if toErr := timeoutCtx.Err(); toErr != nil {
			return toErr
		}
		return errors.New(errMsg)
	case c := <-completedCh:
		if c.Success {
			return nil
		} else {
			return c.Error
		}
	}
}
