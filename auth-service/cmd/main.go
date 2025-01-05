package main

import (
	"github.com/MortalSC/IM-System/auth-service/config"
	"github.com/MortalSC/IM-System/auth-service/internal/router"
	srv "github.com/MortalSC/IM-System/lib/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 路由
	router.InitRouter(r)

	gc := router.RegisterGrpc()
	stop := func() {
		gc.Stop()
	}

	// 启动服务/中止 + grpc服务停止
	srv.RunServer(r, config.Cfg.SrvCfg.Addr, config.Cfg.SrvCfg.Name, stop)
}
