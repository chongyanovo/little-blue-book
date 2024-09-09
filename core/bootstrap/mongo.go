package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

// MongoConfig mongodb配置
type MongoConfig struct {
	Hostname string `mapstructure:"hostname" json:"hostname" yaml:"hostname"` // 服务器地址
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`             // 端口
	Database string `mapstructure:"database" json:"database" yaml:"database"` // 数据库名
	Username string `mapstructure:"username" json:"username" yaml:"username"` // 数据库用户名
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 数据库密码
}

func NewMongo(c *Config, l *zap.Logger) *mongo.Database {
	m := c.MongoConfig
	uri := fmt.Sprintf("mongodb://%v:%v@%v:%v",
		m.Username,
		m.Password,
		m.Hostname,
		m.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			l.Info("MongoDB Command Started",
				zap.String("command", startedEvent.CommandName),
				zap.String("database", startedEvent.DatabaseName),
			)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			l.Info("MongoDB Command Succeeded",
				zap.String("command", succeededEvent.CommandName),
				zap.String("database", succeededEvent.DatabaseName),
			)
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			l.Error("MongoDB Command Failed",
				zap.String("command", failedEvent.CommandName),
				zap.String("database", failedEvent.DatabaseName),
			)
		},
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMonitor(monitor))
	if err != nil {
		panic(err)
	}
	return client.Database(c.MongoConfig.Database)
}
