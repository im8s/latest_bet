package ziface

type IWorkerPool interface {
	Start()
	Stop()

	FetchRequestTo(req IRequest)

	GetRouterHandle() IRouterHandle
}
