package graceful

type Option interface {
	apply(IApplication)
}

type optionFunc func(IApplication)

func (f optionFunc) apply(a IApplication) {
	f(a)
}
