package sms

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type MemoryService struct {
	logger *zap.Logger
}

func NewMemoryService(l *zap.Logger) SmsService {
	return &MemoryService{
		logger: l,
	}
}

func (m *MemoryService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	fmt.Println("====================")
	fmt.Println("验证码：", strings.Join(args, ""))
	fmt.Println("====================")
	return nil
}
