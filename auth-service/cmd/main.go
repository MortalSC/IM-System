package main

import (
	_ "github.com/MortalSC/IM-System/auth-service/api/user"
	"github.com/MortalSC/IM-System/auth-service/config"
	"github.com/MortalSC/IM-System/auth-service/internal/router"
	srv "github.com/MortalSC/IM-System/lib/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 加载配置
	cfg := config.Cfg

	// 路由
	router.InitRouter(r)

	// 启动服务/中止
	srv.RunServer(r, cfg.SrvCfg.Addr, cfg.SrvCfg.Name)
}
