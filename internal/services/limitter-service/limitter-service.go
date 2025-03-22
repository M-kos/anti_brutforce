package limitterservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/M-kos/anti_brutforce/internal/lib/logger"
)

var (
	ErrIpInBlacklist                = "ip is in the blacklist"
	ErrLoginRateLimitterExceeded    = "login rate limitter exceeded"
	ErrPasswordRateLimitterExceeded = "password rate limitter exceeded"
	ErrIpRateLimitterExceeded       = "ip rate limitter exceeded"
	ErrLoginRemoveError             = "login remove error"
	ErrIpRemoveError                = "ip remove error"
)

type RateLimitter interface {
	Check(ctx context.Context, key string) (bool, error)
	Remove(ctx context.Context, key string) error
}

type LabelListChecker interface {
	Check(ctx context.Context, key string) bool
}

type LimitterService struct {
	loginRateLimitter    RateLimitter
	passwordRateLimitter RateLimitter
	ipRateLimitter       RateLimitter

	whitelabelList LabelListChecker
	blacklabelList LabelListChecker

	log logger.LoggerProvider
}

func NewLimitterService(loginRateLimitter, passwordRateLimitter, ipRateLimitter RateLimitter, whitelabelList, blacklabelList LabelListChecker, log logger.LoggerProvider) *LimitterService {
	return &LimitterService{
		loginRateLimitter:    loginRateLimitter,
		passwordRateLimitter: passwordRateLimitter,
		ipRateLimitter:       ipRateLimitter,
		whitelabelList:       whitelabelList,
		blacklabelList:       blacklabelList,
		log:                  log,
	}
}

func (c *LimitterService) Check(ctx context.Context, login string, ip string, password string) error {
	if ok := c.whitelabelList.Check(ctx, ip); ok {
		return nil
	}
	if ok := c.blacklabelList.Check(ctx, ip); ok {
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s: %s", ErrIpInBlacklist, ip))
		return errors.New(ErrIpInBlacklist)
	}

	ok, err := c.loginRateLimitter.Check(ctx, login)
	if err != nil {
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s", err.Error()))
	}

	if !ok {
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s: %s", ErrLoginRateLimitterExceeded, login))
		return errors.New(ErrLoginRateLimitterExceeded)
	}

	ok, err = c.passwordRateLimitter.Check(ctx, password)
	if err != nil {
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s", err.Error()))
	}

	if !ok {
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s: %s", ErrPasswordRateLimitterExceeded, password))
		return errors.New(ErrPasswordRateLimitterExceeded)
	}

	ok, err = c.ipRateLimitter.Check(ctx, ip)
	if err != nil {
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s", err.Error()))
	}

	if !ok {
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s: %s", ErrIpRateLimitterExceeded, ip))
		return errors.New(ErrIpRateLimitterExceeded)
	}

	return nil
}

func (c *LimitterService) Remove(ctx context.Context, login string, ip string) error {
	err := c.loginRateLimitter.Remove(ctx, login)
	if err != nil {
		c.log.Error(fmt.Sprintf("LimitterService: Remove: %s: %s", ErrLoginRemoveError, err.Error()))
	}

	err = c.ipRateLimitter.Remove(ctx, ip)
	if err != nil {
		c.log.Error(fmt.Sprintf("LimitterService: Remove: %s: %s", ErrIpRemoveError, err.Error()))
	}

	return err
}
