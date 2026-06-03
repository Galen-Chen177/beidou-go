package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Login    LoginConfig    `yaml:"login"`
	Channel  ChannelConfig  `yaml:"channel"`
	Script   ScriptConfig   `yaml:"script"`
}

// ServerConfig 通用服务器配置
type ServerConfig struct {
	Name    string `yaml:"name"`    // 服务器名称
	Host    string `yaml:"host"`    // 监听地址
	LogFile string `yaml:"log_file"` // 日志文件路径（空=输出到终端）
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

// DSN 返回 MySQL 连接字符串
func (d DatabaseConfig) DSN() string {
	return d.User + ":" + d.Password +
		"@tcp(" + d.Host + ":" + itoa(d.Port) + ")/" +
		d.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// LoginConfig 登录服务器配置
type LoginConfig struct {
	Port int `yaml:"port"` // 监听端口（默认 8484）
}

// ChannelConfig 频道服务器配置
type ChannelConfig struct {
	Port       int `yaml:"port"`        // 频道端口（默认 7575）
	MaxPlayers int `yaml:"max_players"` // 频道最大玩家数
}

// ScriptConfig 脚本引擎配置
type ScriptConfig struct {
	Path      string `yaml:"path"`       // 脚本目录路径
	HotReload bool   `yaml:"hot_reload"` // 是否启用热加载
}

// Default 返回默认配置
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Name: "BeiDou-Go",
			Host: "0.0.0.0",
		},
		Database: DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Password: "",
			DBName:   "beidou",
		},
		Login: LoginConfig{
			Port: 8484,
		},
		Channel: ChannelConfig{
			Port:       7575,
			MaxPlayers: 100,
		},
		Script: ScriptConfig{
			Path:      "scripts",
			HotReload: false,
		},
	}
}

// Load 从 YAML 文件加载配置。加载失败返回错误。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// MustLoad 加载配置，失败则 panic（用于启动阶段）
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	return cfg
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
