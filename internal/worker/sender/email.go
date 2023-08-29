package sender

import (
	"context"
	"github.com/hookokoko/email"
	"github.com/hookokoko/notify-go/internal"
	"net/smtp"
	"time"
)

type EmailHandler struct {
	EmailClient *email.Client
}

func NewEmailHandler() *EmailHandler {
	ec := email.NewClient(&email.ClientConfig{
		Addr: "smtp.qq.com:25",
		Auth: smtp.PlainAuth("", "648646891@qq.com",
			"", "smtp.qq.com"),
		Options: &email.Options{
			PoolSize:        5,
			PoolTimeout:     30 * time.Second,
			MinIdleConns:    0,
			MaxIdleConns:    1,
			ConnMaxIdleTime: 10 * time.Second, // 距离上一次使用时间多久之后标记失效
		},
	})

	return &EmailHandler{
		EmailClient: ec,
	}
}

func (eh *EmailHandler) Name() string {
	return internal.EmailNAME
}

func (eh *EmailHandler) Execute(ctx context.Context, task *internal.Task) (err error) {
	emailCfg := &email.Email{}
	emailCfg.From = "notifyGo <648646891@qq.com>"
	emailCfg.To = []string{task.MsgReceiver.Value()}
	emailCfg.Text = []byte(task.MsgContent.Content)

	err = eh.EmailClient.SendMail(ctx, emailCfg)

	return
}

//func (eh *EmailHandler) Execute(ctx context.Context, task *internal.Task) (err error) {
//	if task.SendChannel != "email" {
//		return nil
//	}
//  Mock time cost
//	n := tool.RandIntN(700, 800)
//	time.Sleep(time.Millisecond * time.Duration(n))
//	log.Printf("[email] %+v\n, cost: %d ms", task, n)
//	return nil
//}
