package main

import (
	_ "github.com/MortalSC/IM-System/api/user"
	srv "github.com/MortalSC/IM-System/common"
	"github.com/MortalSC/IM-System/config"
	"github.com/MortalSC/IM-System/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 使用 Gin 默认的日志和中间件
	r := gin.Default()

	cfg := config.C

	// 路由
	router.InitRouter(r)

	// 启动服务/终止
	srv.RunServer(r, cfg.SC.Addr, cfg.SC.Name)
}
