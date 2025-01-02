package user

import (
	"github.com/MortalSC/IM-System/router"
	"github.com/gin-gonic/gin"
	"log"
)

// 包初始化函数
func init() {
	log.Println("init user router")
	router.Register(&RouterUser{}) // 注册用户路由
}

// 用户路由结构体
type RouterUser struct {
}

// 实现 Router 接口，将具体路由绑定到 Gin 引擎
func (*RouterUser) Router(r *gin.Engine) {
	h := HandlerUser{}

	// 获取验证码
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
