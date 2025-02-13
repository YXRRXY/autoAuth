package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/YXRRXY/autoAuth/backend/internal/dal/model"
	"github.com/YXRRXY/autoAuth/backend/pkg/utils"
)

type AuthService struct {
	userRepo   UserRepository
	jwtHandler *utils.JWTHandler
}

func NewAuthService(userRepo UserRepository, jwtHandler *utils.JWTHandler) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtHandler: jwtHandler,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

// checkUserExists 检查用户名和邮箱是否已存在
func (s *AuthService) checkUserExists(ctx context.Context, username, email string) error {
	// 检查用户名
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("检查用户名失败: %v", err)
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查邮箱
	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("检查邮箱失败: %v", err)
	}
	if exists {
		return errors.New("邮箱已被注册")
	}

	return nil
}

// Register 处理注册
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) error {
	// 检查用户是否存在
	if err := s.checkUserExists(ctx, req.Username, req.Email); err != nil {
		return err
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   "active",
	}

	return s.userRepo.Create(ctx, user)
}

// Login 处理登录
func (s *AuthService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return "", nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 生成 token
	tokenPair, err := s.jwtHandler.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// 记录日志但不影响登录
		log.Printf("更新最后登录时间失败: %v", err)
	}

	// 清除密码后返回
	user.Password = ""
	return tokenPair.AccessToken, user, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*utils.TokenPair, error) {
	// 使用 JWT 处理器刷新令牌对
	tokenPair, err := s.jwtHandler.RefreshTokenPair(refreshToken)
	if err != nil {
		return nil, err
	}
	return tokenPair, nil
}

// GetUserByID 根据ID获取用户信息
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Password = "" // 清除密码
	return user, nil
}
