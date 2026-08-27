package scripting

import "github.com/dop251/goja"

type GojaPromiseWrapper struct {
	promise *goja.Promise
	resolve func(any) error
	reject  func(any) error
}

func NewGojaPromiseWrapper(promise *goja.Promise, resolve func(any) error, reject func(any) error) *GojaPromiseWrapper {
	return &GojaPromiseWrapper{
		promise: promise,
		resolve: resolve,
		reject:  reject,
	}
}

func isDone(promise *goja.Promise) bool {
	state := promise.State()
	return state == goja.PromiseStateFulfilled || state == goja.PromiseStateRejected
}

func (p *GojaPromiseWrapper) Resolve(value any) error {
	return p.resolve(value)
}

func (p *GojaPromiseWrapper) Reject(value any) error {
	return p.reject(value)
}
