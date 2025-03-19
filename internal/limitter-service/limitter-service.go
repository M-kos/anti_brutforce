package limitterservice

import (
	"context"
	"fmt"
	"log"
)

type Limitter interface {
	Check(ctx context.Context, key string) (bool, error)
	Remove(ctx context.Context, key string) error
}

type LabelList interface {
	Check(ctx context.Context, key string) bool
}

type LimitterService struct {
	loginRateLimitter    Limitter
	passwordRateLimitter Limitter
	ipRateLimitter       Limitter

	whitelabelList LabelList
	blacklabelList LabelList
}

func NewLimitterService(loginRateLimitter, passwordRateLimitter, ipRateLimitter Limitter, whitelabelList, blacklabelList LabelList) *LimitterService {
	return &LimitterService{
		loginRateLimitter:    loginRateLimitter,
		passwordRateLimitter: passwordRateLimitter,
		ipRateLimitter:       ipRateLimitter,
		whitelabelList:       whitelabelList,
		blacklabelList:       blacklabelList,
	}
}

func (c *LimitterService) Check(ctx context.Context, login string, ip string, password string) error {
	if ok := c.whitelabelList.Check(ctx, ip); ok {
		return nil
	}
	if ok := c.blacklabelList.Check(ctx, ip); ok {
		return fmt.Errorf("ip is in the blacklist")
	}

	ok, err := c.loginRateLimitter.Check(ctx, login)
	if err != nil {
		log.Println(err.Error())
	}

	if !ok {
		return fmt.Errorf("login rate limitter exceeded")
	}

	ok, err = c.passwordRateLimitter.Check(ctx, password)
	if err != nil {
		log.Println(err.Error())
	}

	if !ok {
		return fmt.Errorf("password rate limitter exceeded")
	}

	ok, err = c.ipRateLimitter.Check(ctx, ip)
	if err != nil {
		log.Println(err.Error())
	}

	if !ok {
		return fmt.Errorf("ip rate limitter exceeded")
	}

	return nil
}

func (c *LimitterService) Remove(ctx context.Context, login string, ip string) error {
	err := c.loginRateLimitter.Remove(ctx, login)
	if err != nil {
		log.Println("login remove error: %S", err.Error())
	}

	err = c.ipRateLimitter.Remove(ctx, ip)
	if err != nil {
		log.Println("ip remove error: %S", err.Error())
	}

	return err
}
