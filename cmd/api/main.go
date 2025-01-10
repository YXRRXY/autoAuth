package main

import (
	"fmt"
	"log"

	"github.com/YXRRXY/autoAuth/internal/api"
	"github.com/YXRRXY/autoAuth/internal/api/handlers"
	"github.com/YXRRXY/autoAuth/internal/config"
	"github.com/YXRRXY/autoAuth/internal/model"
	"github.com/YXRRXY/autoAuth/internal/repository"
	"github.com/YXRRXY/autoAuth/internal/service"
	"github.com/YXRRXY/autoAuth/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化 JWT
	utils.InitJWTHandler(cfg.JWT.Secret)

	// 初始化各层组件
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, utils.GetJWTHandler())
	authHandler := handlers.NewAuthHandler(authService)

	// 创建服务器
	h := server.Default(server.WithHostPorts(":8080"))

	// 注册路由
	api.RegisterRoutes(h, authHandler)

	// 启动服务器
	h.Spin()
}

// initDB 初始化数据库连接
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 先连接MySQL（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 创建数据库（如果不存在）
	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.Database.DBName)
	if err := db.Exec(createDB).Error; err != nil {
		return nil, err
	}

	// 重新连接（指定数据库）
	dsn = cfg.GetDSN()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
