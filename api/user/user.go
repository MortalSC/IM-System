package user

import (
	"context"
	"fmt"
	common2 "github.com/MortalSC/IM-System/common"
	"github.com/MortalSC/IM-System/common/logs"
	"github.com/MortalSC/IM-System/pkg/dao"
	"github.com/MortalSC/IM-System/pkg/model"
	"github.com/MortalSC/IM-System/pkg/repo"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// HandlerUser 是用户业务处理器结构体
// 封装缓存接口（例如 Redis）以便对数据进行快速存储和检索
type HandlerUser struct {
	cache repo.Cache // 缓存接口，用于存储验证码等临时数据
}

func New() *HandlerUser {
	return &HandlerUser{
		cache: dao.Rc, // 使用全局 RedisCache 实例
	}
}

// getCaptcha 是一个用于获取验证码的 HTTP API 实现
// 它通过手机号码生成验证码，并将其存储到 Redis 中，同时调用短信平台发送验证码【此处模拟实现】
func (h *HandlerUser) getCaptcha(ctx *gin.Context) {
	resp := &common2.Result{}

	// 1. 获取参数
	mobile := ctx.PostForm("mobile") // 解析获取请求表单中的手机号

	// 2. 校验参数
	if !common2.VerifyMobile(mobile) {
		ctx.JSON(http.StatusOK, resp.Failed(model.ErrNoLegalMobile, "手机号不合法"))
		return
	}

	// 3. 生成随机验证码（随机的4位1000~9999或6位100000~999999）
	// 此处为模拟，固定使用 "123456" 作为验证码。
	code := "123456"

	// 4. 调用短信平台
	// 使用 Goroutine 异步调用短信平台，以便快速响应接口请求
	go func() {
		time.Sleep(2 * time.Second) // 模拟调用短信平台的耗时操作
		logs.LG.Info("短信平台调用成功，发送短信")

		// 5. 存储验证码到 Redis
		// 使用带超时的上下文，避免 Redis 请求阻塞主进程
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// 将验证码存储到 Redis，设置过期时间为 15 分钟
		err := h.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			logs.LG.Error(fmt.Sprintf("验证码存入 Redis 出错，原因: %v\n", err))
		}
	}()

	// 6. 返回成功响应（验证码仅供测试时返回，生产环境不应返回验证码）
	ctx.JSON(http.StatusOK, resp.Success(code))
}
