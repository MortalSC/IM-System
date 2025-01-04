package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunServer(r *gin.Engine, addr, serverName string) {
	// 创建一个http.Server实例，使用Gin作为处理器
	// -> 使用标准库http.Server配合Gin，可以更灵活地管理服务（不如自定义超时等）
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在单独的Goroutine中启动服务，使主线程能够继续执行（监听信号、管理关闭逻辑等）
	go func() {
		log.Printf("%s running on port %s\n", serverName, addr)
	}()

	// 捕获关闭信号，用于优雅关闭服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutdown %s ...\n", serverName)

	// 关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown Failed, cause by : %s\n", serverName, err)
	}
	select {
	case <-ctx.Done():
		log.Printf("%s Shutdown timeout\n", serverName)
	}
	log.Printf("%s Shutdown success\n", serverName)
}
