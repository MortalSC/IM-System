package user

import (
	"github.com/MortalSC/IM-System/api-center/internal/router"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	log.Println("init user router")
	ru := &RouterUser{}
	router.Register(ru)
}

// RouterUser ： 用户路由
type RouterUser struct {
}

// Router 方法实现路由注册逻辑
func (ru *RouterUser) Router(r *gin.Engine) {
	// 初始化grpc客户端连接
	InitRpcUserClient()

	h := New()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
