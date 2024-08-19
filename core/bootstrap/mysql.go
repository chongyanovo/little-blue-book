package bootstrap

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

// MysqlConfig mysql配置
type MysqlConfig struct {
	Hostname string `mapstructure:"hostname" json:"hostname" yaml:"hostname"` // 服务器地址
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`             // 端口
	Config   string `mapstructure:"config" json:"config" yaml:"config"`       // 高级配置
	Database string `mapstructure:"database" json:"database" yaml:"database"` // 数据库名
	Username string `mapstructure:"username" json:"username" yaml:"username"` // 数据库用户名
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 数据库密码
}

// NewMysql 创建Mysql数据库连接
func NewMysql(c *Config, l *zap.Logger) *gorm.DB {
	m := c.MysqlConfig
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v",
		m.Username,
		m.Password,
		m.Hostname,
		m.Port,
		m.Database,
		m.Config)
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: dsn,
		}), &gorm.Config{
			Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
				SlowThreshold:             time.Millisecond * 100,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				LogLevel:                  glogger.Info,
			}),
		})
	if err != nil {
		panic(fmt.Sprintf("Mysql创建失败: %v", err))
	}
	return db
}

type gormLoggerFunc func(msg string, args ...zap.Field)

func (g gormLoggerFunc) Printf(msg string, args ...any) {
	g(msg, zap.Field{Key: "args", Interface: args})
}
