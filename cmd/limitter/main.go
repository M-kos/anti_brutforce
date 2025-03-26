package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/M-kos/anti_brutforce/internal/api"
	"github.com/M-kos/anti_brutforce/internal/config"
	"github.com/M-kos/anti_brutforce/internal/db"
	"github.com/M-kos/anti_brutforce/internal/lib/logger"
	"github.com/M-kos/anti_brutforce/internal/lib/ratelimit"
	labellistservice "github.com/M-kos/anti_brutforce/internal/services/label-list-service"
	limitterservice "github.com/M-kos/anti_brutforce/internal/services/limitter-service"
)

func main() {
	go run() // TODO: graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Println("shutting down...")
}

func run() {
	l := logger.NewLogger()
	conf, err := config.LoadConfig()
	if err != nil {
		l.Error("error while creating client: ", err)
		return
	}

	redis := db.NewRedisDB(conf, l)

	loginRateLimitter := ratelimit.NewRateLimiter(conf.LoginNumberAttempts, conf.Timeout, redis)
	passwordRateLimitter := ratelimit.NewRateLimiter(conf.PasswordNumberAttempts, conf.Timeout, redis)
	ipRateLimitter := ratelimit.NewRateLimiter(conf.IPNumberAttempts, conf.Timeout, redis)

	whitelabelList := labellistservice.NewLabelListService(
		labellistservice.NewLabelListRepository(redis, labellistservice.WhitelistKey),
		l,
		labellistservice.WhitelistKey,
	)
	blacklabelList := labellistservice.NewLabelListService(
		labellistservice.NewLabelListRepository(redis, labellistservice.BlacklistKey),
		l,
		labellistservice.BlacklistKey,
	)

	limitterService := limitterservice.NewLimitterService(
		loginRateLimitter,
		passwordRateLimitter,
		ipRateLimitter,
		whitelabelList,
		blacklabelList,
		l,
	)

	err = api.Run(conf, limitterService, whitelabelList, blacklabelList, l)
	if err != nil {
		l.Error(err.Error())
	}
}
