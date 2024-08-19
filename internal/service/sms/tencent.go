package sms

import (
	"context"
	"errors"
	"fmt"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"math/rand"
)

type TencentSmsService struct {
	appId     string
	signature string
	client    *sms.Client
}

func NewTencentSmsService(client *sms.Client, appId string, signature string) SmsService {
	return &TencentSmsService{
		appId:     appId,
		signature: signature,
		client:    client,
	}
}

func (s TencentSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = &s.appId
	req.SignName = &s.signature
	req.TemplateId = &tplId
	req.TemplateParamSet = s.toStringPtrSlice(args)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *status.Code != "Ok" {
			return errors.New(
				fmt.Sprintf("发送短信失败，%s，%s", *status.Code, *status.Message),
			)
		}
	}
	return nil
}

func (s TencentSmsService) toStringPtrSlice(slice []string) []*string {
	res := make([]*string, len(slice))
	for i, s := range slice {
		res[i] = &s
	}
	return res
}

func (s TencentSmsService) generateCode() string {
	num := rand.Intn(999999)
	return fmt.Sprintf("%6d", num)
}
