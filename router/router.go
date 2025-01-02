package router

import (
	"github.com/gin-gonic/gin"
)

// 定义接口，用于抽象路由行为
type Router interface {
	Router(r *gin.Engine)
}

// 路由注册器的结构体
type RegisterRouter struct {
}

// 工厂函数，用于创建路由注册器实例
func New() *RegisterRouter {
	return &RegisterRouter{}
}

// 路由注册方法
func (*RegisterRouter) Router(ro Router, r *gin.Engine) {
	ro.Router(r) // 调用具体路由的 Router 方法
}

// 存储所有注册的路由
var routers []Router

//// 【写法一】
//func InitRouter(r *gin.Engine) {
//	register := New()
//	register.Router(&user.RouterUser{}, r)
//}

// 【写法二】初始化路由，将所有注册的路由绑定到 Gin 引擎
func InitRouter(r *gin.Engine) {
	for _, ro := range routers {
		ro.Router(r)
	}
}

// 注册路由（可变参数）
func Register(ro ...Router) {
	routers = append(routers, ro...)
}
