package bootstrap

type CacheConfig struct {
	UserExpiration int64 `mapstructure:"user-expiration" json:"user-expiration" yaml:"user-expiration"`
}
