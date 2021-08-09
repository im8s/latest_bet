package ziface

type IRouterHandle interface {
	DoHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	MatchRouter(msgID uint32) bool

	// FetchRequestTo(request IRequest)
}
