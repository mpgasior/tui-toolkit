package async

import "context"

type Context interface {
	context.Context
	Id() FuncId
	YieldEvent(e any)
	YieldFunc(id FuncId, fn Func)
}

type contextImpl struct {
	context.Context
	cancel context.CancelFunc
	runner *Runner
	id     FuncId
}

func (c *contextImpl) Id() FuncId {
	return c.id
}

func (c *contextImpl) YieldEvent(e any) {
	c.runner.yieldEvent(e)
}

func (c *contextImpl) YieldFunc(id FuncId, fn Func) {
	c.runner.Run(id, fn)
}
