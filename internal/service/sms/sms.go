package sms

import "context"

type SmsService interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
