package sender

import (
	"context"
	"github.com/hookokoko/notify-go/internal"
	"github.com/hookokoko/notify-go/pkg/tool"
	"log"
	"time"
)

type SmsHandler struct{}

func NewSmsHandler() *SmsHandler {
	return &SmsHandler{}
}

func (fh *SmsHandler) Name() string {
	return internal.SmsNAME
}

func (fh *SmsHandler) Execute(ctx context.Context, task *internal.Task) error {
	if task.SendChannel != "sms" {
		return nil
	}
	// Mock time cost
	n := tool.RandIntN(700, 800)
	time.Sleep(time.Millisecond * time.Duration(n))
	log.Printf("[sms] %+v\n, cost: %d ms", task, n)
	return nil
}
