package blacklist

import "context"

type Repository interface {
	GetList(ctx context.Context, listKey string) ([]string, error)
	AddToList(ctx context.Context, listKey string, value string) error
	RemoveFromList(ctx context.Context, listKey string, value string) error
}

const (
	blacklistKey = "blacklist"
)

type BlacklistRepo struct {
	repo Repository
}

func NewBlacklistRepository(repository Repository) *BlacklistRepo {
	return &BlacklistRepo{
		repo: repository,
	}
}

func (wr *BlacklistRepo) Get(ctx context.Context) ([]string, error) {
	values, err := wr.repo.GetList(ctx, blacklistKey)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (wr *BlacklistRepo) Add(ctx context.Context, value string) error {
	err := wr.repo.AddToList(ctx, blacklistKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (wr *BlacklistRepo) Remove(ctx context.Context, value string) error {
	err := wr.repo.RemoveFromList(ctx, blacklistKey, value)
	if err != nil {
		return err
	}

	return nil
}
