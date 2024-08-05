package memory

import (
	"context"
	"fmt"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	fmt.Println("====================")
	fmt.Println("验证码：", strings.Join(args, ""))
	fmt.Println("====================")
	return nil
}
