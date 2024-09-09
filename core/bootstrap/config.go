package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 配置文件
type Config struct {
	ServerConfig *ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	ZapConfig    *ZapConfig    `mapstructure:"zap" json:"zap" yaml:"zap"`
	MysqlConfig  *MysqlConfig  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	MongoConfig  *MongoConfig  `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
	RedisConfig  *RedisConfig  `mapstructure:"redis" json:"redis" yaml:"redis"`
	TokenConfig  *TokenConfig  `mapstructure:"token" json:"token" yaml:"token"`
	CacheConfig  *CacheConfig  `mapstructure:"cache" json:"cache" yaml:"cache"`
	LimitConfig  *LimitConfig  `mapstructure:"limit" json:"limit" yaml:"limit"`
}

// NewConfig 读取配置文件
func NewConfig(v *viper.Viper) *Config {
	c := Config{}
	if err := v.Unmarshal(&c); err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}
	return &c
}
