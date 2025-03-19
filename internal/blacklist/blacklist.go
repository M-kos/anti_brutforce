package blacklist

import (
	"context"
	"fmt"
	"log"
	"net"
)

type BlacklistRepository interface {
	Get(ctx context.Context) ([]string, error)
	Add(ctx context.Context, value string) error
	Remove(ctx context.Context, value string) error
}

type Blacklist struct {
	repo BlacklistRepository
}

func NewBlacklist(repository BlacklistRepository) *Blacklist {
	return &Blacklist{
		repo: repository,
	}
}

func (b *Blacklist) Check(ctx context.Context, value string) bool {
	cidrs, err := b.repo.Get(ctx)
	if err != nil {
		log.Printf("Blacklist: check wrong %s", err)
	}

	for _, v := range cidrs {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			log.Printf("Blacklist: check wrong %s", err)
			continue
		}

		if ipNet.Contains(net.ParseIP(value)) {
			return true
		}
	}

	return false
}

func (b *Blacklist) Add(ctx context.Context, value string) error {
	if err := b.repo.Add(ctx, value); err != nil {
		log.Printf("Blacklist: value not added: %s", err.Error())
		return fmt.Errorf("not added")
	}

	return nil
}

func (b *Blacklist) Remove(ctx context.Context, value string) error {
	if err := b.repo.Remove(ctx, value); err != nil {
		log.Printf("Blacklist: value not removed: %s", err.Error())
		return fmt.Errorf("not removed")
	}

	return nil
}
