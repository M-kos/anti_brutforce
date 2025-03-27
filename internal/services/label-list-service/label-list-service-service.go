package labellistservice

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/M-kos/anti_brutforce/internal/models"
)

var (
	ErrCheck  = "Check: check wrong"
	ErrAdd    = "Add: value not added"
	ErrRemove = "Remove: value not removed"
)

const (
	WhitelistKey = "whitelist"
	BlacklistKey = "blacklist"
)

type LabelListRepositoryProvider interface {
	Get(ctx context.Context) ([]string, error)
	Add(ctx context.Context, value string) error
	Remove(ctx context.Context, value string) error
}

type LabelListService struct {
	repository LabelListRepositoryProvider
	log        models.LoggerProvider
	key        string
}

func NewLabelListService(
	repository LabelListRepositoryProvider,
	log models.LoggerProvider,
	key string,
) *LabelListService {
	return &LabelListService{
		repository: repository,
		log:        log,
		key:        key,
	}
}

func (ll *LabelListService) Check(ctx context.Context, value string) bool {
	cidrs, err := ll.repository.Get(ctx)
	if err != nil {
		ll.log.Error(fmt.Sprintf("%s: %s: %s", ll.key, ErrCheck, err.Error()))
	}

	for _, v := range cidrs {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			ll.log.Error(fmt.Sprintf("%s: %s: %s", ll.key, ErrCheck, err.Error()))
			continue
		}

		if ipNet.Contains(net.ParseIP(value)) {
			return true
		}
	}

	return false
}

func (ll *LabelListService) Add(ctx context.Context, value string) error {
	if err := ll.repository.Add(ctx, value); err != nil {
		ll.log.Error(fmt.Sprintf("%s: %s: %s", ll.key, ErrAdd, err.Error()))
		return errors.New(ErrAdd)
	}

	return nil
}

func (ll *LabelListService) Remove(ctx context.Context, value string) error {
	if err := ll.repository.Remove(ctx, value); err != nil {
		ll.log.Error(fmt.Sprintf("%s: %s: %s", ll.key, ErrRemove, err.Error()))
		return errors.New(ErrRemove)
	}

	return nil
}
