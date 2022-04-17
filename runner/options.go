package runner

import (
	"os"
	"time"
)

type Option interface {
	apply(*funcOption)
	signal() chan os.Signal
	task() chan Task
	timeouts() <-chan time.Time
}

type funcOption struct {
	f  func(*funcOption)
	to <-chan time.Time
	s  chan os.Signal
	t  chan Task
}

func (fo *funcOption) task() chan Task {
	return fo.t
}

func (fo *funcOption) timeouts() <-chan time.Time {
	return fo.to
}

func (fo *funcOption) signal() chan os.Signal {
	return fo.s
}

func (fo *funcOption) apply(do *funcOption) {
	fo.f(do)
}

func newFuncOption(f func(*funcOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func parseOptions(opts ...Option) Option {
	o := new(funcOption)
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func WithTimeout(t <-chan time.Time) Option {
	return newFuncOption(func(options *funcOption) {
		options.to = t
	})
}

func WithSingle(s chan os.Signal) Option {
	return newFuncOption(func(options *funcOption) {
		options.s = s
	})
}

func WithTask(t chan Task) Option {
	return newFuncOption(func(options *funcOption) {
		options.t = t
	})
}
