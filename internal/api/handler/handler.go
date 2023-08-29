package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hookokoko/notify-go/internal"
	"github.com/hookokoko/notify-go/internal/config"
	"github.com/hookokoko/notify-go/internal/service"
	"github.com/hookokoko/notify-go/pkg/mq"

	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	core *service.Core
}

func NewPushHandler(mqCfg *mq.Config, dbCfg config.DBConfig) *PushHandler {
	c := service.NewCore(mqCfg, dbCfg)
	return &PushHandler{
		core: c,
	}
}

// SendBatch TODO 批量发送消息待实现
func (p *PushHandler) SendBatch(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "send batch")
}

func (p *PushHandler) Send(ctx *gin.Context) {
	params := make(map[string]interface{}, 16)
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "绑定参数出错")
	}

	targets := make([]internal.ITarget, 0, 16)

	channel := params["channel"].(string)
	switch channel {
	case "email":
		emails, ok := params["email"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, "获取邮件地址失败")
			return
		}
		for _, email := range strings.Split(emails, ",") {
			targets = append(targets, internal.EmailTarget{Email: email})
		}
	case "sms":
		phones, ok := params["phone"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, "获取手机号失败")
			return
		}
		for _, phone := range strings.Split(phones, ",") {
			targets = append(targets, internal.PhoneTarget{Phone: phone})
		}
	case "push":
		userIds, ok := params["userId"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, "获取用户id失败")
			return
		}
		for _, uid := range strings.Split(userIds, ",") {
			targets = append(targets, internal.IdTarget{Id: uid})
		}
	default:
		ctx.JSON(http.StatusBadRequest, "不支持的发送渠道")
		return
	}

	templateIdStr := params["templateId"].(string)
	templateId, _ := strconv.ParseInt(templateIdStr, 10, 64)

	err = p.core.Send(ctx, channel, targets, templateId, map[string]interface{}{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "消息处理处理失败")
		return
	}
	ctx.JSON(http.StatusOK, "send ok")

	return
}
