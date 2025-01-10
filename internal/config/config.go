package config

import (
	"fmt"
	"os"

	"github.com/YXRRXY/autoAuth/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host" default:"localhost"`
		Port     int    `yaml:"port" default:"3306"`
		Username string `yaml:"username" default:"root"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname" default:"blocksbom"`
	} `yaml:"database"`

	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// 读取配置文件
	configPath := "../../configs/config.yaml"
	if err := loadYamlConfig(configPath, cfg); err != nil {
		return nil, err
	}

	// 如果JWT密钥为空，生成新密钥并更新配置文件
	if cfg.JWT.Secret == "" {
		secret, err := utils.GenerateRandomSecret()
		if err != nil {
			return nil, fmt.Errorf("generate secret failed: %v", err)
		}
		cfg.JWT.Secret = secret

		// 将更新后的配置写回文件
		newData, err := yaml.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("marshal config error: %v", err)
		}

		if err := os.WriteFile(configPath, newData, 0644); err != nil {
			return nil, fmt.Errorf("update config file error: %v", err)
		}
	}

	return cfg, nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
	)
}

// loadYamlConfig 从YAML文件加载配置
func loadYamlConfig(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}
