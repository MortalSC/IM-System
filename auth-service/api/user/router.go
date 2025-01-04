package user

import (
	"github.com/MortalSC/IM-System/auth-service/internal/router"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	log.Println("User Router initialized")
	router.Register(&RouterUser{})
}

// RouterUser ： 用户路由
type RouterUser struct {
}

func (*RouterUser) Router(r *gin.Engine) {
	h := &HandlerUser{}
	log.Println("Registering /project/login/getCaptcha route") // 添加调试日志
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
