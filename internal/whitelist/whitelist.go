package whitelist

import (
	"context"
	"fmt"
	"log"
	"net"
)

type WhitelistRepository interface {
	Get(ctx context.Context) ([]string, error)
	Add(ctx context.Context, value string) error
	Remove(ctx context.Context, value string) error
}

type Whitelist struct {
	repository WhitelistRepository
}

func NewWhitelist(repository WhitelistRepository) *Whitelist {
	return &Whitelist{
		repository: repository,
	}
}

func (w *Whitelist) Check(ctx context.Context, value string) bool {
	cidrs, err := w.repository.Get(ctx)
	if err != nil {
		log.Printf("Whitelist: check wrong %s", err)
	}

	for _, v := range cidrs {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			log.Printf("Whitelist: check wrong %s", err)
			continue
		}

		if ipNet.Contains(net.ParseIP(value)) {
			return true
		}
	}

	return false
}

func (w *Whitelist) Add(ctx context.Context, value string) error {
	if err := w.repository.Add(ctx, value); err != nil {
		log.Printf("Whitelist: value not added: %s", err.Error())
		return fmt.Errorf("not added")
	}

	return nil
}

func (w *Whitelist) Remove(ctx context.Context, value string) error {
	if err := w.repository.Remove(ctx, value); err != nil {
		log.Printf("Whitelist: value not removed: %s", err.Error())
		return fmt.Errorf("not removed")
	}

	return nil
}
