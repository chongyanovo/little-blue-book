package bootstrap

import "time"

type LimitConfig struct {
	SmsLimitConfig *SmsLimitConfig `mapstructure:"sms" json:"sms" yaml:"sms"`
}

type SmsLimitConfig struct {
	Interval time.Duration `mapstructure:"interval" json:"interval" yaml:"interval"`
	Rate     int           `mapstructure:"rate" json:"rate" yaml:"rate"`
}
