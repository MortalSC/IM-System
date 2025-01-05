package login_service_v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/MortalSC/IM-System/auth-service/internal/utils"
	LibCache "github.com/MortalSC/IM-System/lib/cache"
	libLog "github.com/MortalSC/IM-System/lib/log"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	cache LibCache.Cache
}

func New(cache LibCache.Cache) *LoginService {
	return &LoginService{
		cache: cache,
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *CaptchaMessage) (*CaptchaResponse, error) {
	// 1. 获取参数
	mobile := msg.Mobile // 解析获取请求表单中的手机号

	// 2. 校验参数
	if !utils.VerifyMobile(mobile) {
		return nil, errors.New("手机号不合法")
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
		err := ls.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			libLog.IMLog.Error(fmt.Sprintf("验证码存入 Redis 出错，原因: %v\n", err))
		}
	}()

	return &CaptchaResponse{Code: code}, nil
}
