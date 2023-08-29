package sender

import (
	"context"
	"github.com/hookokoko/notify-go/internal"
	"github.com/hookokoko/notify-go/pkg/item"
)

type IHandler interface {
	Name() string
	Execute(ctx context.Context, taskInfo *internal.Task) error
}

type HandleManager struct {
	manager *item.Manager
}

func NewHandlerManager() *HandleManager {
	return &HandleManager{
		manager: item.NewManager(
			NewEmailHandler(),
			NewSmsHandler(),
			NewPushHandler(),
		),
	}
}

func (hm *HandleManager) Get(key string) (resp IHandler, err error) {
	if h, err := hm.manager.Get(key); err == nil {
		return h.(IHandler), nil
	}
	return nil, err
}
