package blacklist

type Repository interface {
	Get(key string) (string, error)
	Add(key string) error
	Remove(key string) error
}

type BlacklistRepo struct {
	repo Repository
}

func NewBlacklistRepository(repository Repository) *BlacklistRepo {
	return &BlacklistRepo{
		repo: repository,
	}
}

func (br *BlacklistRepo) Get(key string) (string, error) {
	value, err := br.repo.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (br *BlacklistRepo) Add(key string) error {
	err := br.repo.Add(key)
	if err != nil {
		return err
	}

	return nil
}

func (br *BlacklistRepo) Remove(key string) error {
	err := br.repo.Remove(key)
	if err != nil {
		return err
	}

	return nil
}
