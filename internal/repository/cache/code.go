package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrSetCodeTooMany         = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknowErrCode          = errors.New("我也不知道发生什么，反正是和code有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里面
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	redis  redis.Cmdable
	logger *zap.Logger
}

func NewCodeCache(r redis.Cmdable, l *zap.Logger) *CodeCache {
	return &CodeCache{
		redis:  r,
		logger: l,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.redis.Eval(ctx, luaSetCode, []string{c.generateKey(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
	case -1:
		// 发送太频繁
		return ErrSetCodeTooMany
	default:
		// 系统错误
		return errors.New("系统错误")
	}
	return err
}

func (c *CodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.redis.Eval(ctx, luaVerifyCode, []string{c.generateKey(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，你需要告警，因为有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	default:
		return false, ErrUnknowErrCode
	}
}

func (c *CodeCache) generateKey(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
