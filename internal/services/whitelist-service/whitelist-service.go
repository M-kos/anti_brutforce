package whitelistservice

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/M-kos/anti_brutforce/internal/lib/logger"
)

var (
	ErrCheck  = "Whitelist: Check: check wrong"
	ErrAdd    = "Whitelist: Add: value not added"
	ErrRemove = "Whitelist: Remove: value not removed"
)

type WhitelistRepositoryProvider interface {
	Get(ctx context.Context) ([]string, error)
	Add(ctx context.Context, value string) error
	Remove(ctx context.Context, value string) error
}

type WhitelistService struct {
	repository WhitelistRepositoryProvider
	log        logger.LoggerProvider
}

func NewWhitelist(repository WhitelistRepositoryProvider, log logger.LoggerProvider) *WhitelistService {
	return &WhitelistService{
		repository: repository,
		log:        log,
	}
}

func (w *WhitelistService) Check(ctx context.Context, value string) bool {
	cidrs, err := w.repository.Get(ctx)
	if err != nil {
		w.log.Error(fmt.Sprintf("%s: %s", ErrCheck, err.Error()))
	}

	for _, v := range cidrs {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			w.log.Error(fmt.Sprintf("%s: %s", ErrCheck, err.Error()))
			continue
		}

		if ipNet.Contains(net.ParseIP(value)) {
			return true
		}
	}

	return false
}

func (w *WhitelistService) Add(ctx context.Context, value string) error {
	if err := w.repository.Add(ctx, value); err != nil {
		w.log.Error(fmt.Sprintf("%s: %s", ErrAdd, err.Error()))
		return errors.New(ErrAdd)
	}

	return nil
}

func (w *WhitelistService) Remove(ctx context.Context, value string) error {
	if err := w.repository.Remove(ctx, value); err != nil {
		w.log.Error(fmt.Sprintf("%s: %s", ErrRemove, err.Error()))
		return errors.New(ErrRemove)
	}

	return nil
}
