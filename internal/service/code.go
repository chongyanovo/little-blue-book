package service

import (
	"context"
	"fmt"
	"github.com/ChongYanOvO/little-blue-book/internal/repository"
	"github.com/ChongYanOvO/little-blue-book/internal/service/sms"
	"go.uber.org/zap"
	"math/rand"
)

var codeTpId = "1877556"

type CodeService interface {
	Send(ctx context.Context,
		biz string,
		phone string) error
	Verify(ctx context.Context, biz string,
		phone string, inputCode string) (bool, error)
}

type CodeServiceImpl struct {
	repo   repository.CodeRepository
	smsSvc sms.SmsService
	logger *zap.Logger
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.SmsService, l *zap.Logger) CodeService {
	return &CodeServiceImpl{
		repo:   repo,
		smsSvc: smsSvc,
		logger: l,
	}
}

// Send 发验证码 我需要什么参数
func (svc *CodeServiceImpl) Send(ctx context.Context,
	// 区别使用业务
	biz string,
	phone string) error {
	// 生成一个验证码
	code := svc.generateCode()
	// 塞进去 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送出去
	err = svc.smsSvc.Send(ctx, codeTpId, []string{code}, phone)
	if err != nil {
		// 这意味着，Redis 有这个验证码，但是不好意思
		// 我能不能删掉这个验证码
		// 你这个 err 可能是超时的 err，你都不知道收到了吗
		// 在这里重试，要重试的时候，传入一个自己就会重试的 smsSvcRetry

	}
	return err
}

func (svc *CodeServiceImpl) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeServiceImpl) generateCode() string {
	num := rand.Intn(999999)
	// 不够 6 位的，加上前导 0
	return fmt.Sprintf("%6d", num)
}
