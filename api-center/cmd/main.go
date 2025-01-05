package main

import (
	_ "github.com/MortalSC/IM-System/api-center/api/user"
	"github.com/MortalSC/IM-System/api-center/config"
	"github.com/MortalSC/IM-System/api-center/internal/router"
	srv "github.com/MortalSC/IM-System/lib/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 路由
	router.InitRouter(r)

	// 启动服务/中止 + grpc服务停止
	srv.RunServer(r, config.Cfg.SrvCfg.Addr, config.Cfg.SrvCfg.Name, nil)
}
