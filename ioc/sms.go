package ioc

import (
	"github.com/ChongYanOvO/little-blue-book/internal/service/sms"
	"github.com/ChongYanOvO/little-blue-book/internal/service/sms/memory"
)

func InitSmsService() sms.Service {
	return memory.NewService()
}
