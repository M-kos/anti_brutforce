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
	blacklistservice "github.com/M-kos/anti_brutforce/internal/services/blacklist-service"
	limitterservice "github.com/M-kos/anti_brutforce/internal/services/limitter-service"
	whitelistservice "github.com/M-kos/anti_brutforce/internal/services/whitelist-service"
)

func main() {
	go run() // TODO: graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Println("shutting down...")
}

func run() error {
	conf := config.LoadConfig()
	l := logger.NewLogger()
	redis := db.NewRedisDb(conf, l)

	loginRateLimitter := ratelimit.NewRateLimiter(conf.LoginNumberAttempts, conf.Timeout, redis)
	passwordRateLimitter := ratelimit.NewRateLimiter(conf.PasswordNumberAttempts, conf.Timeout, redis)
	ipRateLimitter := ratelimit.NewRateLimiter(conf.IPNumberAttempts, conf.Timeout, redis)

	whitelabelList := whitelistservice.NewWhitelist(whitelistservice.NewWhitelistRepository(redis), l)
	blacklabelList := blacklistservice.NewBlacklist(blacklistservice.NewBlacklistRepository(redis), l)

	limitterService := limitterservice.NewLimitterService(loginRateLimitter, passwordRateLimitter, ipRateLimitter, whitelabelList, blacklabelList, l)

	err := api.Run(conf, limitterService, whitelabelList, blacklabelList, l)
	if err != nil {
		l.Error(err.Error())
	}

	return err
}
