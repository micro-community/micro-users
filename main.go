package main

import (
	"github.com/micro-community/micro-users/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	_ "github.com/micro-community/micro-users/profile"
)

func main() {
	srv := service.New(
		service.Name("micro-users"),
	)

	srv.Init(service.BeforeStart(func() error {
		redisHostValue, err := config.Get("Redis.Host")
		redisHost := redisHostValue.String("")

		if err != nil && redisHost != "" {
			logger.Info("config:", redisHost)
		}

		return nil
	}))

	srv.Handle(handler.NewUsers())

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
