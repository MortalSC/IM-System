package user

import (
	"github.com/MortalSC/IM-System/auth-service/config"
	"github.com/MortalSC/IM-System/auth-service/internal/router"
	"github.com/MortalSC/IM-System/lib/cache"
	"github.com/MortalSC/IM-System/lib/cache/redis"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	cfg := config.Cfg
	// 初始化 Redis 缓存实例
	cacheInstance, err := redis.NewRedisCache(cfg.InitRedisOptions())
	if err != nil {
		log.Fatalf("Failed to initialize Redis cache: %v\n", err)
	}

	log.Println("User Router initialized")
	// 注册路由，同时注入依赖
	router.Register(NewRouterUser(cacheInstance))
}

// RouterUser ： 用户路由
type RouterUser struct {
	cache cache.Cache // 引入缓存接口作为依赖
}

// NewRouterUser 创建一个 RouterUser 实例
// 参数：
// - cache: Cache 接口实现，用于依赖注入
// 返回值：
// - *RouterUser: RouterUser 的实例
func NewRouterUser(cache cache.Cache) *RouterUser {
	return &RouterUser{cache: cache}
}

// Router 方法实现路由注册逻辑
func (ru *RouterUser) Router(r *gin.Engine) {
	// 将缓存实例注入到 HandlerUser 中
	h := NewHandlerUser(ru.cache)
	log.Println("Registering /project/login/getCaptcha route") // 添加调试日志
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
