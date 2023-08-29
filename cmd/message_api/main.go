package main

import (
	"log"
	"net/http"

	"github.com/hookokoko/notify-go/internal/api/router"
	"github.com/hookokoko/notify-go/internal/config"
	"github.com/hookokoko/notify-go/pkg/mq"
)

func main() {
	addr := ":8080"
	srv := newServer(addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("启动失败{%v}", err)
		}
		log.Fatalf("服务启动, 端口：%s...", addr)
	}()

	// TODO 服务优雅退出
}

func newServer(addr string) *http.Server {
	mqCfg := mq.NewConfig("./conf/kafka_topic.toml")
	appCfg := config.LoadConfig("./conf/app.toml")

	pusher := router.NewMsgPusher(mqCfg, appCfg.DB["default"])
	return &http.Server{
		Addr:    addr,
		Handler: pusher.GetRouter(),
	}
}
