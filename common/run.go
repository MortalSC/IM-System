package common

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

	// 创建一个http.Server实例，使用 Gin 作为处理器
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	// 使用标准库 http.Server 配合 Gin，可以更灵活地管理服务（比如自定义超时等）。

	// 在单独的 Goroutine 中启动服务，使主线程能够继续执行（监听信号、管理关闭逻辑等）
	go func() {
		log.Printf("%s running on port %s\n", serverName, srv.Addr)
		// 捕获 ListenAndServe 的错误，除了 http.ErrServerClosed（这是正常关闭的标志），其他错误会被记录并导致程序终止。
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 使用 Go 的 os.Signal 和 signal.Notify 监听 SIGINT（用户按下 Ctrl+C）和 SIGTERM（一般用于容器关闭信号）。
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting Down project web server...")

	// 关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %s\n", serverName, err)
	}
	select {
	case <-ctx.Done():
		log.Printf("%s Shutdown timeout\n", serverName)
	}
	log.Printf("%s stop success...\n", serverName)

}
