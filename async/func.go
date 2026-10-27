package async

type FuncId string

type Func func(ctx Context) (any, error)

type FuncCompletedEvent struct {
	Id   FuncId
	Data any
	Err  error
}
