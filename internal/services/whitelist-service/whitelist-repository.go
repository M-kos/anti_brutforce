package whitelistservice

import (
	"context"

	"github.com/M-kos/anti_brutforce/internal/models"
)

type DbProvider interface {
	GetList(ctx context.Context, listKey string) (*models.LabelList, error)
	AddToList(ctx context.Context, listKey string, value string) error
	RemoveFromList(ctx context.Context, listKey string, value string) error
}

const (
	WhitelistKey = "whitelist"
)

type WhitelistRepo struct {
	db DbProvider
}

func NewWhitelistRepository(repository DbProvider) *WhitelistRepo {
	return &WhitelistRepo{
		db: repository,
	}
}

func (wr *WhitelistRepo) Get(ctx context.Context) ([]string, error) {
	list, err := wr.db.GetList(ctx, WhitelistKey)
	if err != nil {
		return nil, err
	}

	return list.Values, nil
}

func (wr *WhitelistRepo) Add(ctx context.Context, value string) error {
	err := wr.db.AddToList(ctx, WhitelistKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (wr *WhitelistRepo) Remove(ctx context.Context, value string) error {
	err := wr.db.RemoveFromList(ctx, WhitelistKey, value)
	if err != nil {
		return err
	}

	return nil
}
