package blacklistservice

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/M-kos/anti_brutforce/internal/models"
)

var (
	ErrCheck  = "blacklist: Check: check wrong"
	ErrAdd    = "blacklist: Add: value not added"
	ErrRemove = "blacklist: Remove: value not removed"
)

type BlacklistRepositoryProvider interface {
	Get(ctx context.Context) ([]string, error)
	Add(ctx context.Context, value string) error
	Remove(ctx context.Context, value string) error
}

type BlacklistService struct {
	repository BlacklistRepositoryProvider
	log        models.LoggerProvider
}

func NewBlacklist(repository BlacklistRepositoryProvider, log models.LoggerProvider) *BlacklistService {
	return &BlacklistService{
		repository: repository,
		log:        log,
	}
}

func (b *BlacklistService) Check(ctx context.Context, value string) bool {
	cidrs, err := b.repository.Get(ctx)
	if err != nil {
		b.log.Error(fmt.Sprintf("%s: %s", ErrCheck, err.Error()))
	}

	for _, v := range cidrs {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			b.log.Error(fmt.Sprintf("%s: %s", ErrCheck, err.Error()))
			continue
		}

		if ipNet.Contains(net.ParseIP(value)) {
			return true
		}
	}

	return false
}

func (b *BlacklistService) Add(ctx context.Context, value string) error {
	if err := b.repository.Add(ctx, value); err != nil {
		b.log.Error(fmt.Sprintf("%s: %s", ErrAdd, err.Error()))
		return errors.New(ErrAdd)
	}

	return nil
}

func (b *BlacklistService) Remove(ctx context.Context, value string) error {
	if err := b.repository.Remove(ctx, value); err != nil {
		b.log.Error(fmt.Sprintf("%s: %s", ErrRemove, err.Error()))
		return errors.New(ErrRemove)
	}

	return nil
}
