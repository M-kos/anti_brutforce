package whitelist

type Repository interface {
	Get(key string) (string, error)
	Add(key string) error
	Remove(key string) error
}

type WhitelistRepo struct {
	repo Repository
}

func NewWhitelistRepository(repository Repository) *WhitelistRepo {
	return &WhitelistRepo{
		repo: repository,
	}
}

func (wr *WhitelistRepo) Get(key string) (string, error) {
	value, err := wr.repo.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (wr *WhitelistRepo) Add(key string) error {
	err := wr.repo.Add(key)
	if err != nil {
		return err
	}

	return nil
}

func (wr *WhitelistRepo) Remove(key string) error {
	err := wr.repo.Remove(key)
	if err != nil {
		return err
	}

	return nil
}
