package main

import "fmt"

func main() {
	app, err := InitApp()
	if err != nil {
		panic(err)
	}

	config, err := InitConfig()
	if err != nil {
		panic(err)
	}
	app.Server.Run(
		fmt.Sprintf("%s:%d",
			config.ServerConfig.Host,
			config.ServerConfig.Port),
	)
}
