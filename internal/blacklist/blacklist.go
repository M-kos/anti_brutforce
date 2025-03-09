package blacklist

type Store interface {
	Check(key string) (bool, error)
	Add(key string) (bool, error)
	Remove(key string) (bool, error)
}

type Blacklist struct {
	store Store
}

func NewBlacklist(store Store) *Blacklist {
	return &Blacklist{
		store: store,
	}
}

func (b *Blacklist) Check(key string) (bool, error) {
	return b.store.Check(key)
}

func (b *Blacklist) Add(key string) (bool, error) {
	return b.store.Add(key)
}

func (b *Blacklist) Remove(key string) (bool, error) {
	return b.store.Remove(key)
}
