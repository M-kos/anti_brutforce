package api

import "errors"

var (
	ErrLoginIsRequired     = errors.New("login is required")
	ErrPasswordIsRequired  = errors.New("password is required")
	ErrIpIsRequired        = errors.New("ip is required")
	ErrCidrRequired        = errors.New("cidr is required")
	ErrCheck               = errors.New("check error")
	ErrRemove              = errors.New("remove error")
	ErrAddToWhitelist      = errors.New("err adding to whitelist error")
	ErrRemoveFromWhitelist = errors.New("err removing from whitelist error")
	ErrAddToBlacklist      = errors.New("err adding to blacklist error")
	ErrRemoveFromBlacklist = errors.New("err removing from blacklist error")
)
