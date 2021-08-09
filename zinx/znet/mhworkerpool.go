package znet

import (
	"context"
	"fmt"

	"zinx/ziface"
)

type Worker struct {
	WP       ziface.IWorkerPool
	JobQueue chan ziface.IRequest
}

func NewWorker(wp ziface.IWorkerPool) (wk *Worker) {
	return &Worker{
		WP:       wp,
		JobQueue: make(chan ziface.IRequest),
	}
}

func (wk *Worker) Run(wq chan chan ziface.IRequest, ctx context.Context) {
	go func() {
		for {
			wq <- wk.JobQueue

			select {
			case req, ok := <-wk.JobQueue:
				if ok {
					// fmt.Println("[before] wk.WP.GetRouterHandle().DoHandler: msgID = ", req.GetMsgID())
					wk.WP.GetRouterHandle().DoHandler(req)
					// fmt.Println("[after] wk.WP.GetRouterHandle().DoHandler: msgID = ", req.GetMsgID())
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

type WorkerPool struct {
	Workerlen   uint32
	JobQueue    chan ziface.IRequest
	WorkerQueue chan chan ziface.IRequest

	MH ziface.IRouterHandle

	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorkerPool(workerlen uint32) (wp *WorkerPool) {
	wp = &WorkerPool{
		Workerlen:   workerlen,
		JobQueue:    make(chan ziface.IRequest),
		WorkerQueue: make(chan chan ziface.IRequest, workerlen),
		MH:          NewRouterHandle(),
	}

	return wp
}

func (wp *WorkerPool) GetRouterHandle() ziface.IRouterHandle {
	return wp.MH
}

func (wp *WorkerPool) Start() {
	fmt.Println("init workerpool")

	wp.ctx, wp.cancel = context.WithCancel(context.Background())

	for i := 0; i < int(wp.Workerlen); i++ {
		worker := NewWorker(wp)
		worker.Run(wp.WorkerQueue, wp.ctx)
	}

	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.WorkerQueue
				worker <- job
			}
		}
	}()
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
}

func (wp *WorkerPool) FetchRequestTo(req ziface.IRequest) {
	wp.JobQueue <- req
}
