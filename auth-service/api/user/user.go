package user

import (
	"context"
	"fmt"
	"github.com/MortalSC/IM-System/auth-service/internal/errors"
	"github.com/MortalSC/IM-System/auth-service/internal/utils"
	"github.com/MortalSC/IM-System/auth-service/pkg/model"
	LibCache "github.com/MortalSC/IM-System/lib/cache"
	libLog "github.com/MortalSC/IM-System/lib/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// HandlerUser 是用户业务处理器结构体
type HandlerUser struct {
	cache LibCache.Cache
}

// NewHandlerUser 创建一个新的 HandlerUser 实例
// 参数：
// - cache: Cache 接口实现（可以是 Redis 实现，也可以是其他实现）
// 返回值：
// - *HandlerUser: HandlerUser 的实例
func NewHandlerUser(cache LibCache.Cache) *HandlerUser {
	return &HandlerUser{
		cache: cache,
	}
}

// getCaptcha 是一个用于获取验证码的 HTTP API 实现
// [POST] /project/login/getCaptcha
func (h *HandlerUser) getCaptcha(ctx *gin.Context) {
	resp := model.HttpResult{}

	// 1. 获取参数
	mobile := ctx.PostForm("mobile") // 解析获取请求表单中的手机号

	// 2. 校验参数
	if !utils.VerifyMobile(mobile) {
		ctx.JSON(http.StatusOK, resp.Failed(errors.ErrNoLegalMobile, "手机号不合法"))
		return
	}

	// 3. 生成随机验证码（随机的4位1000~9999或6位100000~999999）
	// 此处为模拟，固定使用 "123456" 作为验证码。
	code := "123456"

	// 4. 调用短信平台
	// 使用 Goroutine 异步调用短信平台，以便快速响应接口请求
	go func() {
		time.Sleep(2 * time.Second) // 模拟调用短信平台的耗时操作
		libLog.IMLog.Info("短信平台调用成功，发送短信")

		// 5. 存储验证码到 Redis
		// 使用带超时的上下文，避免 Redis 请求阻塞主进程
		c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// 将验证码存储到 Redis，设置过期时间为 15 分钟
		err := h.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			libLog.IMLog.Error(fmt.Sprintf("验证码存入 Redis 出错，原因: %v\n", err))
		}
	}()

	// 6. 返回成功响应（验证码仅供测试时返回，生产环境不应返回验证码）
	ctx.JSON(http.StatusOK, resp.Success(code))
}
