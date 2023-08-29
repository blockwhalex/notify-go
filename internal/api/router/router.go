package router

import (
	"github.com/hookokoko/notify-go/internal/api/handler"
	"github.com/hookokoko/notify-go/internal/config"
	"github.com/hookokoko/notify-go/pkg/mq"

	"github.com/gin-gonic/gin"
)

type MsgPusher struct {
	PushHandler *handler.PushHandler
}

func NewMsgPusher(mqCfg *mq.Config, dbCfg config.DBConfig) *MsgPusher {
	return &MsgPusher{PushHandler: handler.NewPushHandler(mqCfg, dbCfg)}
}

func (m *MsgPusher) GetRouter() *gin.Engine {
	router := gin.New()
	g := router.Group("/message")

	// TODO 接口鉴权

	// 发送消息
	g.POST("send", m.PushHandler.Send)
	g.POST("sendBatch", m.PushHandler.SendBatch)

	// TODO 查看消息记录
	//g.GET("Send")
	//g.GET("SendBatch")

	return router
}
