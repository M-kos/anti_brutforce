package whitelist

import "context"

type Repository interface {
	GetList(ctx context.Context, listKey string) ([]string, error)
	AddToList(ctx context.Context, listKey string, value string) error
	RemoveFromList(ctx context.Context, listKey string, value string) error
}

const (
	whitelistKey = "whitelist"
)

type WhitelistRepo struct {
	repo Repository
}

func NewWhitelistRepository(repository Repository) *WhitelistRepo {
	return &WhitelistRepo{
		repo: repository,
	}
}

func (wr *WhitelistRepo) Get(ctx context.Context) ([]string, error) {
	values, err := wr.repo.GetList(ctx, whitelistKey)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (wr *WhitelistRepo) Add(ctx context.Context, value string) error {
	err := wr.repo.AddToList(ctx, whitelistKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (wr *WhitelistRepo) Remove(ctx context.Context, value string) error {
	err := wr.repo.RemoveFromList(ctx, whitelistKey, value)
	if err != nil {
		return err
	}

	return nil
}
