package blacklistservice

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
	BlacklistKey = "blacklist"
)

type BlacklistRepository struct {
	db DbProvider
}

func NewBlacklistRepository(db DbProvider) *BlacklistRepository {
	return &BlacklistRepository{
		db: db,
	}
}

func (wr *BlacklistRepository) Get(ctx context.Context) ([]string, error) {
	list, err := wr.db.GetList(ctx, BlacklistKey)
	if err != nil {
		return nil, err
	}

	return list.Values, nil
}

func (wr *BlacklistRepository) Add(ctx context.Context, value string) error {
	err := wr.db.AddToList(ctx, BlacklistKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (wr *BlacklistRepository) Remove(ctx context.Context, value string) error {
	err := wr.db.RemoveFromList(ctx, BlacklistKey, value)
	if err != nil {
		return err
	}

	return nil
}
