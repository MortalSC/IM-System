package user

import (
	"fmt"
	"github.com/MortalSC/IM-System/api-center/pkg/model"
	loginServiceV1 "github.com/MortalSC/IM-System/auth-service/pkg/service/login.service.v1"
	libLog "github.com/MortalSC/IM-System/lib/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandlerUser 是用户业务处理器结构体
type HandlerUser struct {
}

func New() *HandlerUser {
	return &HandlerUser{}
}

// getCaptcha 是一个用于获取验证码的 HTTP API 实现
// [POST] /project/login/getCaptcha
func (h *HandlerUser) getCaptcha(ctx *gin.Context) {
	result := model.HttpResult{}

	mobile := ctx.PostForm("mobile")

	if mobile == "" {
		ctx.JSON(http.StatusOK, result.Failed(2002, "手机号不能为空"))
		return
	}
	libLog.IMLog.Info(fmt.Sprintf("请求验证码，手机号：%s", mobile))

	_, err := LoginServiceClient.GetCaptcha(ctx, &loginServiceV1.CaptchaMessage{Mobile: mobile})
	if err != nil {
		libLog.IMLog.Error(fmt.Sprintf("调用 GetCaptcha 出错：%v", err))
		ctx.JSON(http.StatusOK, result.Failed(2001, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(result.Code))

}
