package router

import (
	"fmt"
	"github.com/MortalSC/IM-System/auth-service/config"
	loginServiceV1 "github.com/MortalSC/IM-System/auth-service/pkg/service/login.service.v1"
	"github.com/MortalSC/IM-System/lib/cache"
	"github.com/MortalSC/IM-System/lib/cache/redis"
	libLog "github.com/MortalSC/IM-System/lib/log"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
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

func RegisterGrpc() *grpc.Server {
	// 初始化依赖
	cacheInstance, err := InitDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	c := gRPCConfig{
		Addr: config.Cfg.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			loginServiceV1.RegisterLoginServiceServer(g, loginServiceV1.New(cacheInstance))
		},
	}

	s := grpc.NewServer()
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}

// InitDependencies 初始化服务依赖
func InitDependencies() (cache.Cache, error) {
	cfg := config.Cfg
	cacheInstance, err := redis.NewRedisCache(cfg.InitRedisOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}
	libLog.IMLog.Debug("Redis initialized successfully")
	return cacheInstance, nil
}
