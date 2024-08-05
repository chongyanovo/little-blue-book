package config

var Config = config{
	DB: DBConfig{
		DSN: "root:123456@tcp(localhost:13306)/little-blue-book",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
