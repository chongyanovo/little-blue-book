package bootstrap

import (
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Address  string `mapstructure:"address" json:"address" yaml:"address"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

func NewRedis(c *Config) redis.Cmdable {
	r := c.RedisConfig
	return redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Username: r.User,
		Password: r.Password,
	})
}
