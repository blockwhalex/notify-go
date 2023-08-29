package service

import (
	"context"
	"github.com/hookokoko/notify-go/internal"
)

type TargetService struct{}

type ITargetService interface {
	GetTarget(ctx context.Context, targetId int64) []internal.ITarget
}

func NewTargetService() ITargetService {
	return &TargetService{}
}

func (ts *TargetService) GetTarget(ctx context.Context, targetId int64) []internal.ITarget {
	targets := []internal.ITarget{
		internal.IdTarget{Id: "111"},
		internal.EmailTarget{Email: "ch_hakun@163.com"},
		internal.PhoneTarget{Phone: "+8618800187099"},
	}
	return targets
}
