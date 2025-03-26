package labellistservice

import (
	"context"

	"github.com/M-kos/anti_brutforce/internal/models"
)

type DBProvider interface {
	GetList(ctx context.Context, listKey string) (*models.LabelList, error)
	AddToList(ctx context.Context, listKey string, value string) error
	RemoveFromList(ctx context.Context, listKey string, value string) error
}

type LabelListRepository struct {
	db  DBProvider
	key string
}

func NewLabelListRepository(repository DBProvider, key string) *LabelListRepository {
	return &LabelListRepository{
		db:  repository,
		key: key,
	}
}

func (ll *LabelListRepository) Get(ctx context.Context) ([]string, error) {
	list, err := ll.db.GetList(ctx, ll.key)
	if err != nil {
		return nil, err
	}

	return list.Values, nil
}

func (ll *LabelListRepository) Add(ctx context.Context, value string) error {
	err := ll.db.AddToList(ctx, ll.key, value)
	if err != nil {
		return err
	}

	return nil
}

func (ll *LabelListRepository) Remove(ctx context.Context, value string) error {
	err := ll.db.RemoveFromList(ctx, ll.key, value)
	if err != nil {
		return err
	}

	return nil
}
