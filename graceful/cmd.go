package graceful

import (
	"context"
	"time"
)

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

func WithName(name string) Option {
	return optionFunc(func(a IApplication) {
		a.SetName(name)
	})
}

func WithCmd(cmd string) Option {
	return optionFunc(func(a IApplication) {
		a.SetCmd(cmd)
	})
}

func WithWorkPath(wPath string) Option {
	return optionFunc(func(a IApplication) {
		a.SetWorkPath(wPath)
	})
}

func WithTimeOut(timeout time.Duration) Option {
	return optionFunc(func(a IApplication) {
		a.SetTimeOut(timeout)
	})
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(a IApplication) {
		a.SetContext(ctx)
	})
}
