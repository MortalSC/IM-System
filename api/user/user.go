package user

import (
	common2 "github.com/MortalSC/IM-System/common"
	"github.com/MortalSC/IM-System/pkg/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// 用户业务处理器
type HandlerUser struct {
}

// 一个简单的 API 实现
func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	resp := &common2.Result{}

	// 1. 获取参数
	mobile := ctx.PostForm("mobile") // 解析获取请求表单中的手机号

	// 2. 校验参数
	if !common2.VerifyMobile(mobile) {
		ctx.JSON(http.StatusOK, resp.Failed(model.ErrNoLegalMobile, "手机号不合法"))
		return
	}

	// 3. 生成随机验证码（随机的4位1000~9999或6位100000~999999）
	code := "123456"

	// 4. 调用短信平台（三方，放入go协程中执行，接口可以快速响应）
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("短信平台调用成功，发送短信")
		// 5. 存储验证码（放在redis中，过期时间为15分钟）
		log.Printf("将手机号和验证码存入redis成功： PEGISTER_%s : %s", mobile, code)
	}()

	ctx.JSON(http.StatusOK, resp.Success(code))
}
