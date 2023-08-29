package sender

import (
	"context"
	"github.com/hookokoko/notify-go/internal"
	"github.com/hookokoko/notify-go/pkg/tool"
	"log"
	"time"
)

type PushHandler struct{}

func NewPushHandler() *PushHandler {
	return &PushHandler{}
}

func (fh *PushHandler) Name() string {
	return internal.PushNAME
}

func (fh *PushHandler) Execute(ctx context.Context, task *internal.Task) error {
	if task.SendChannel != "push" {
		return nil
	}
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	log.Printf("[push] %+v\n, cost: %d ms", task, n)
	return nil
}
