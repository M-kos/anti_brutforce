package stubrepo

type Repository interface {
	Get(key string) (string, error)
	Add(key string) error
	Remove(key string) error
}

type StubRepo struct {
	list []string
}

func NewStubRepo() Repository {
	return &StubRepo{}
}

func (s *StubRepo) Get(key string) (string, error) {
	for _, v := range s.list {
		if v == key {
			return v, nil
		}
	}
	return "", nil
}

func (s *StubRepo) Add(key string) error {
	s.list = append(s.list, key)
	return nil
}

func (s *StubRepo) Remove(key string) error {
	for i, v := range s.list {
		if v == key {
			s.list = append(s.list[:i], s.list[i+1:]...)
			return nil
		}
	}
	return nil
}
