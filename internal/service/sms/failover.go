package sms

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"sync/atomic"
)

type FailOverService struct {
	logger *zap.Logger
	svcs   []SmsService
	index  uint64
}

func NewFailOverService(logger *zap.Logger, svcs []SmsService) SmsService {
	return FailOverService{
		svcs:   svcs,
		logger: logger,
	}
}

func (f FailOverService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.index, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < length; i++ {
		err := f.svcs[i%length].Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			f.logger.Info("发送短信成功")
			return nil
		case context.DeadlineExceeded, context.Canceled:
			f.logger.Error("发送短信失败", zap.Error(err))
			return err
		default:
			f.logger.Error("发送短信异常", zap.Error(err))
		}
	}
	return errors.New("短信服务全部失败")
}
