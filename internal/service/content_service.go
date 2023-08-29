package service

import (
	"context"

	"github.com/hookokoko/notify-go/internal"
	"github.com/hookokoko/notify-go/internal/config"
	"github.com/hookokoko/notify-go/internal/model"
)

/*
content服务 根据模版组装发送内容
*/

type ContentService struct {
	TemplateDAO model.ITemplateDAO
}

type IContentService interface {
	GetContent(ctx context.Context, target internal.ITarget,
		templateId int64, variable map[string]interface{}) internal.MsgContent
}

func NewContentService(dbCfg config.DBConfig) IContentService {
	return &ContentService{
		TemplateDAO: model.NewITemplateDAO(dbCfg),
	}
}

func (cs *ContentService) GetContent(ctx context.Context, target internal.ITarget,
	templateId int64,
	variable map[string]interface{}) internal.MsgContent {
	content := internal.MsgContent{}
	msg, err := cs.TemplateDAO.GetContent(templateId, "")
	if err != nil {
		return content
	}
	content.Content = msg
	return content
}
