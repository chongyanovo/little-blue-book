package main

import (
	"fmt"
	"github.com/ChongYanOvO/little-blue-book/wire"
)

func main() {
	app, err := wire.InitApp()
	if err != nil {
		panic(err)
	}

	config, err := wire.InitConfig()
	if err != nil {
		panic(err)
	}
	app.Server.Run(
		fmt.Sprintf("%s:%d",
			config.ServerConfig.Host,
			config.ServerConfig.Port),
	)
}
