package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Router : 定义接口，用于抽象路由行为
type Router interface {
	Router(r *gin.Engine)
}

// RouterRegister : 路由注册器
type RouterRegister struct {
}

// New : 工厂函数，用于创建路由注册器实例
func New() *RouterRegister {
	return &RouterRegister{}
}

// Router : 路由注册方法
func (*RouterRegister) Router(ro Router, r *gin.Engine) {
	ro.Router(r)
}

// 存储所有注册的路由
var routers []Router

// InitRouter : 初始化路由，将所有的路由绑定到Gin引擎
func InitRouter(r *gin.Engine) {
	for _, ro := range routers {
		ro.Router(r)
	}
}

// Register : 路由注册（可变参数）
func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}
