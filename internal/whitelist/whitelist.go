package whitelist

import (
	"context"
	"fmt"
	"log"
)

type WhitelistRepository interface {
	Get(key string) (string, error)
	Add(key string) error
	Remove(key string) error
}

type Whitelist struct {
	repository WhitelistRepository
}

func NewWhitelist(repository WhitelistRepository) *Whitelist {
	return &Whitelist{
		repository: repository,
	}
}

func (w *Whitelist) Check(key string) bool {
	value, err := w.repository.Get(key)
	if err != nil {
		log.Printf("Whitelist: check wrong %s", err)
	}

	return value != ""
}

func (w *Whitelist) Add(ctx context.Context, key string) error {
	if err := w.repository.Add(key); err != nil {
		log.Printf("Whitelist: key not added: %s", err.Error())
		return fmt.Errorf("not added")
	}

	return nil
}

func (w *Whitelist) Remove(ctx context.Context, key string) error {
	if err := w.repository.Remove(key); err != nil {
		log.Printf("Whitelist: key not removed: %s", err.Error())
		return fmt.Errorf("not removed")
	}

	return nil
}
