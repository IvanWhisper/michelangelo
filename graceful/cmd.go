package graceful

import (
	"context"
	"time"
)

// NewIApplication create an application
func NewIApplication(opts ...Option) IApplication {
	app := &application{}
	for _, o := range opts {
		o.apply(app)
	}
	if app.ctx == nil {
		app.SetContext(context.Background())
	}
	return app
}

// WithName named application
func WithName(name string) Option {
	return optionFunc(func(a IApplication) {
		a.SetName(name)
	})
}

// WithCmd out app call command
func WithCmd(cmd string) Option {
	return optionFunc(func(a IApplication) {
		a.SetCmd(cmd)
	})
}

// WithPrintCh use output chan, if nil then use std
func WithPrintCh() Option {
	return optionFunc(func(a IApplication) {
		ch := make(chan string)
		a.SetPrintCh(ch)
	})
}

// WithWorkPath work dir path
func WithWorkPath(wPath string) Option {
	return optionFunc(func(a IApplication) {
		a.SetWorkPath(wPath)
	})
}

// WithTimeOut work must in set time
func WithTimeOut(timeout time.Duration) Option {
	return optionFunc(func(a IApplication) {
		a.SetTimeOut(timeout)
	})
}

// WithContext run context
func WithContext(ctx context.Context) Option {
	return optionFunc(func(a IApplication) {
		a.SetContext(ctx)
	})
}
