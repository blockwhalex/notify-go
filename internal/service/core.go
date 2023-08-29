package service

import (
	"context"
	"fmt"

	"github.com/hookokoko/notify-go/internal"
	"github.com/hookokoko/notify-go/internal/config"
	"github.com/hookokoko/notify-go/pkg/mq"

	"golang.org/x/sync/errgroup"
)

// 请求target服务，获取发送目标
// 请求content服务，获取发送内容
// 发送到mq

type Core struct {
	ContentService IContentService
	TargetService  ITargetService
	SendService    IService
}

func NewCore(mqCfg *mq.KafkaConfig, dbCfg config.DBConfig) *Core {
	return &Core{
		ContentService: NewContentService(dbCfg),
		TargetService:  NewTargetService(),
		SendService:    NewSendService(mqCfg, dbCfg),
	}
}

// Send
// 1. 创建一个delivery记录, 同时会获取所有的target，针对每一个target创建一条target记录关联delivery id
// 2. 推送至kafka
func (c *Core) Send(ctx context.Context, channel string, targets []internal.ITarget,
	templateId int64, variable map[string]interface{}) error {
	batchTask := make([]internal.Task, 0, len(targets))

	var eg errgroup.Group
	for _, recvr := range targets {
		eg.Go(func() error {
			msgContent := c.ContentService.GetContent(ctx, recvr, templateId, variable)
			batchTask = append(batchTask, internal.Task{
				//MsgId:       0,
				SendChannel: channel,
				MsgContent:  msgContent,
				MsgReceiver: recvr,
			})
			return c.SendService.Process(ctx, templateId, batchTask)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("get error:%v\n", err)
	}
	return nil
}

// SendBatch 的区别是利用target服务获取待发送目标，如用户画像
// 批量的话是不是可以一边获取target一边发送，即流式推送
func (c *Core) SendBatch(ctx context.Context, channel string, targetId,
	templateId int64, variable map[string]interface{}) error {

	// 通过目标服务服务发送目标，比如用户画像
	// TODO 上百万级别的发送目标方案考虑
	receivers := c.TargetService.GetTarget(ctx, targetId)

	batchTask := make([]internal.Task, 0, len(receivers))

	var eg errgroup.Group
	for _, recvr := range receivers {
		// TODO 针对每一个发送目标构建发送内容，后续需要考虑加一个缓存
		// TODO 当发送目标过多，需要使用goroutine池限制住goroutine
		eg.Go(func() error {
			msgContent := c.ContentService.GetContent(ctx, recvr, templateId, variable)
			batchTask = append(batchTask, internal.Task{
				//MsgId:       0,
				SendChannel: channel,
				MsgContent:  msgContent,
				MsgReceiver: recvr,
			})
			return c.SendService.Process(ctx, templateId, batchTask)
		})
	}
	if err := eg.Wait(); err != nil {
		fmt.Printf("get error:%v\n", err)
	}
	return nil
}
