package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/M-kos/anti_brutforce/internal/api"
	"github.com/M-kos/anti_brutforce/internal/blacklist"
	bucketlist "github.com/M-kos/anti_brutforce/internal/bucket"
	"github.com/M-kos/anti_brutforce/internal/config"
	"github.com/M-kos/anti_brutforce/internal/db"
	limitterservice "github.com/M-kos/anti_brutforce/internal/limitter-service"
	"github.com/M-kos/anti_brutforce/internal/whitelist"
	"github.com/M-kos/anti_brutforce/pkg/ratelimit"
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

	redis := db.NewDb()

	bucketList := bucketlist.NewDbBucketList(redis)

	loginRateLimitter := ratelimit.NewRateLimiter(conf.LoginNumberAttempts, conf.Timeout, bucketList)
	passwordRateLimitter := ratelimit.NewRateLimiter(conf.PasswordNumberAttempts, conf.Timeout, bucketList)
	ipRateLimitter := ratelimit.NewRateLimiter(conf.IPNumberAttempts, conf.Timeout, bucketList)

	whitelabelList := whitelist.NewWhitelist(whitelist.NewWhitelistRepository(bucketList))
	blacklabelList := blacklist.NewBlacklist(blacklist.NewBlacklistRepository(bucketList))

	limitterService := limitterservice.NewLimitterService(loginRateLimitter, passwordRateLimitter, ipRateLimitter, whitelabelList, blacklabelList)

	err := api.Run(conf, limitterService, whitelabelList, blacklabelList)

	return err
}
