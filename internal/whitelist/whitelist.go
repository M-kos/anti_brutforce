package whitelist

type Store interface {
	Check(key string) (bool, error)
	Add(key string) (bool, error)
	Remove(key string) (bool, error)
}

type Whitelist struct {
	store Store
}

func NewWhitelist(store Store) *Whitelist {
	return &Whitelist{
		store: store,
	}
}

func (w *Whitelist) Check(key string) (bool, error) {
	return w.store.Check(key)
}

func (w *Whitelist) Add(key string) (bool, error) {
	return w.store.Add(key)
}

func (w *Whitelist) Remove(key string) (bool, error) {
	return w.store.Remove(key)
}
