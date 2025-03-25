package limitterservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/M-kos/anti_brutforce/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIPInBlacklist                = "ip is in the blacklist"
	ErrLoginRateLimitterExceeded    = "login rate limitter exceeded"
	ErrPasswordRateLimitterExceeded = "password rate limitter exceeded"
	ErrIPRateLimitterExceeded       = "ip rate limitter exceeded"
	ErrLoginRemoveError             = "login remove error"
	ErrIPRemoveError                = "ip remove error"
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

	log models.LoggerProvider
}

func NewLimitterService(
	loginRateLimitter,
	passwordRateLimitter,
	ipRateLimitter RateLimitter,
	whitelabelList,
	blacklabelList LabelListChecker,
	log models.LoggerProvider,
) *LimitterService {
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
		c.log.Error(fmt.Sprintf("LimitterService: Check: %s: %s", ErrIPInBlacklist, ip))
		return errors.New(ErrIPInBlacklist)
	}

	ok, err := c.loginRateLimitter.Check(ctx, login)
	if err != nil {
		c.log.Error(fmt.Sprintf("limitterService: Check: %s", err.Error()))
	}

	if !ok {
		c.log.Error(fmt.Sprintf("limitterService: Check: %s: %s", ErrLoginRateLimitterExceeded, login))
		return errors.New(ErrLoginRateLimitterExceeded)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		pass = []byte(password)
	}

	ok, err = c.passwordRateLimitter.Check(ctx, string(pass))
	if err != nil {
		c.log.Error(fmt.Sprintf("limitterService: Check: %s", err.Error()))
	}

	if !ok {
		c.log.Error(fmt.Sprintf("limitterService: Check: %s: %s (hash: %s)", ErrPasswordRateLimitterExceeded, password, pass))
		return errors.New(ErrPasswordRateLimitterExceeded)
	}

	ok, err = c.ipRateLimitter.Check(ctx, ip)
	if err != nil {
		c.log.Error(fmt.Sprintf("limitterService: Check: %s", err.Error()))
	}

	if !ok {
		c.log.Error(fmt.Sprintf("limitterService: Check: %s: %s", ErrIPRateLimitterExceeded, ip))
		return errors.New(ErrIPRateLimitterExceeded)
	}

	return nil
}

func (c *LimitterService) Remove(ctx context.Context, login string, ip string) error {
	err := c.loginRateLimitter.Remove(ctx, login)
	if err != nil {
		c.log.Error(fmt.Sprintf("limitterService: Remove: %s: %s", ErrLoginRemoveError, err.Error()))
	}

	err = c.ipRateLimitter.Remove(ctx, ip)
	if err != nil {
		c.log.Error(fmt.Sprintf("limitterService: Remove: %s: %s", ErrIPRemoveError, err.Error()))
	}

	return err
}
