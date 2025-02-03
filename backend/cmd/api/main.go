package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YXRRXY/autoAuth/backend/internal/api"
	"github.com/YXRRXY/autoAuth/backend/internal/api/handlers"
	"github.com/YXRRXY/autoAuth/backend/internal/config"
	"github.com/YXRRXY/autoAuth/backend/internal/dal"
	"github.com/YXRRXY/autoAuth/backend/internal/dal/query"
	"github.com/YXRRXY/autoAuth/backend/internal/service"
	"github.com/YXRRXY/autoAuth/backend/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func setupApp(cfg *config.Config) (*server.Hertz, error) {
	// 初始化数据库
	if err := dal.InitDB(cfg); err != nil {
		return nil, err
	}

	// 初始化 JWT
	utils.InitJWTHandler(cfg.JWT.Secret)

	// 初始化各层组件
	userRepo := query.NewUserRepository(dal.GetDB())
	authService := service.NewAuthService(userRepo, utils.GetJWTHandler())
	authHandler := handlers.NewAuthHandler(authService)

	// 创建服务器
	h := server.Default(server.WithHostPorts(":8080"))

	// 注册路由
	api.RegisterRoutes(h, authHandler)

	return h, nil
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化应用
	h, err := setupApp(cfg)
	if err != nil {
		log.Fatalf("Failed to setup application: %v", err)
	}
	defer dal.Close()

	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理信号
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("正在关闭服务器...")
		
		// 创建一个带超时的上下文
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()
		
		if err := h.Shutdown(timeoutCtx); err != nil {
			log.Printf("服务器关闭出错: %v\n", err)
		}
		cancel()
	}()

	// 启动服务器
	go h.Spin()

	// 等待关闭信号
	<-ctx.Done()
	log.Println("服务器已关闭")
}
