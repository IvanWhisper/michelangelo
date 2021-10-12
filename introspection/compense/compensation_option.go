package compense

import (
	"context"
	"github.com/IvanWhisper/michelangelo/dependencies"
	"time"

	"xorm.io/xorm"
)

type Option interface {
	apply(compensation *Compensation)
}

// type optionFunc func
/**
 * @Description:
 * @param compensation
 */
type optionFunc func(compensation *Compensation)

// apply
/**
 * @Description:
 * @receiver f
 * @param c
 */
func (f optionFunc) apply(c *Compensation) {
	f(c)
}

// WithDB
/**
 * @Description:
 * @param v
 * @return Option
 */
func WithDB(v *xorm.Engine) Option {
	return optionFunc(func(compensation *Compensation) {
		compensation.db = v
	})
}

// WithMargin
/**
 * @Description:
 * @param v
 * @return Option
 */
func WithMargin(v time.Duration) Option {
	return optionFunc(func(compensation *Compensation) {
		compensation.margin = v
	})
}

// WithFuncMap
/**
 * @Description:
 * @param v
 * @return Option
 */
func WithFuncMap(v map[string]dependencies.Exec) Option {
	return optionFunc(func(compensation *Compensation) {
		compensation.funcMap = v
	})
}

// WithExecCallBack
/**
 * @Description:
 * @param v
 * @return Option
 */
func WithExecCallBack(v func(string)) Option {
	return optionFunc(func(compensation *Compensation) {
		compensation.execCallBack = v
	})
}

// WithNewAsyncCtx
/**
 * @Description:
 * @param v
 * @return Option
 */
func WithNewAsyncCtx(v func(ctx context.Context) context.Context) Option {
	return optionFunc(func(compensation *Compensation) {
		compensation.newAsyncCtx = v
	})
}

// NewCompensation
/**
 * @Description:
 * @param opts
 * @return dependencies.ICompensator
 */
func NewCompensation(opts ...Option) dependencies.ICompensator {
	c := &Compensation{}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}
