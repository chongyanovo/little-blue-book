package sms

import (
	"context"
	"fmt"
	"github.com/ChongYanOvO/little-blue-book/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var errLimit = fmt.Errorf("触发限流")

type LimiterService struct {
	svc     SmsService
	limiter ratelimit.Limiter
	logger  *zap.Logger
}

func NewSmsLimiterService(cmd redis.Cmdable, svc SmsService) ratelimit.Limiter {
	return &ratelimit.RedisSlidingWindowLimiter{
		Cmd:      cmd,
		Interval: 10 * 60 * 1000,
		Rate:     10,
	}
}

func NewLimiterService(svc SmsService, limiter ratelimit.Limiter, logger *zap.Logger) SmsService {
	return &LimiterService{
		svc:     svc,
		limiter: limiter,
		logger:  logger,
	}

}

func (l LimiterService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limit, err := l.limiter.Limit(ctx, "sms")
	if err != nil {
		l.logger.Error("短信服务限流出现问题", zap.Error(err))
		return err
	}
	if limit {
		l.logger.Warn("触发限流")
		return errLimit
	}
	return l.svc.Send(ctx, tplId, args, numbers...)
}
