package blacklist

import (
	"context"
	"fmt"
	"log"
)

type BlacklistRepository interface {
	Get(key string) (string, error)
	Add(key string) error
	Remove(key string) error
}

type Blacklist struct {
	repo BlacklistRepository
}

func NewBlacklist(repository BlacklistRepository) *Blacklist {
	return &Blacklist{
		repo: repository,
	}
}

func (b *Blacklist) Check(key string) bool {
	value, err := b.repo.Get(key)
	if err != nil {
		log.Printf("Blacklist: check wrong %s", err)
	}

	return value != ""
}

func (b *Blacklist) Add(ctx context.Context, key string) error {
	if err := b.repo.Add(key); err != nil {
		log.Printf("Blacklist: key not added: %s", err.Error())
		return fmt.Errorf("not added")
	}

	return nil
}

func (b *Blacklist) Remove(ctx context.Context, key string) error {
	if err := b.repo.Remove(key); err != nil {
		log.Printf("Blacklist: key not removed: %s", err.Error())
		return fmt.Errorf("not removed")
	}

	return nil
}
