package store

type storageRepository interface {
	Get(key string) (string, error)
	Add(key string) error
	Remove(key string) error
}

type Store struct {
	repository storageRepository
}

func NewStore(repository storageRepository) *Store {
	return &Store{
		repository: repository,
	}
}

func (s *Store) Get(key string) (string, error) {
	value, err := s.repository.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *Store) Add(key string) (string, error) {
	err := s.repository.Add(key)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *Store) Remove(key string) (string, error) {
	err := s.repository.Remove(key)
	if err != nil {
		return "", err
	}

	return key, nil
}
