package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hookokoko/notify-go/internal"
	"github.com/hookokoko/notify-go/internal/config"
	"github.com/hookokoko/notify-go/internal/model"
	"github.com/hookokoko/notify-go/pkg/mq"
)

type IService interface {
	Process(ctx context.Context, templateId int64, tasks []internal.Task) error
}

type SendService struct {
	producer    *mq.Producer
	NotifyGoDAO model.INotifyGoDAO
}

func NewSendService(mqCfg *mq.KafkaConfig, dbCfg config.DBConfig) IService {
	return &SendService{
		producer:    mq.NewProducer(mqCfg),
		NotifyGoDAO: model.NewINotifyGoDAO(dbCfg),
	}
}

func (ss *SendService) Process(ctx context.Context, templateId int64, tasks []internal.Task) error {
	for _, task := range tasks {
		//
		err := ss.NotifyGoDAO.InsertRecord(ctx, templateId, task.MsgReceiver, task.MsgContent.Content)
		if err != nil {
			log.Println("insert send record ")
			continue
		}

		taskBytes, err := json.Marshal(task)
		if err != nil {
			return err
		}
		ss.producer.Send(task.SendChannel, taskBytes)
	}

	return nil
}
