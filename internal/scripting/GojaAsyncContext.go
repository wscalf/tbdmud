package scripting

type GojaAsyncContext struct {
	requeue func()
	promise *GojaPromiseWrapper
	result  any
	err     error
}

func NewGojaAsyncContext(requeue func()) *GojaAsyncContext {
	return &GojaAsyncContext{
		requeue: requeue,
	}
}

func (c *GojaAsyncContext) SetPromise(promise *GojaPromiseWrapper) {
	c.promise = promise
}

func (c *GojaAsyncContext) SetResult(result any) {
	c.result = result
	c.requeue()
}

func (c *GojaAsyncContext) SetError(err error) {
	c.err = err
	c.requeue()
}

func (c *GojaAsyncContext) Resolve() error {
	if c.err != nil {
		return c.promise.Reject(c.err)
	} else {
		return c.promise.Resolve(c.result)
	}
}
