package bucketlist

import "time"

type InMemoryucketList struct {
	store map[string]struct {
		count     uint
		startTime time.Time
	}
}

func NewInMemoryucket() *InMemoryucketList {
	return &InMemoryucketList{
		store: make(map[string]struct {
			count     uint
			startTime time.Time
		}),
	}
}

func (b *InMemoryucketList) Set(key string, count uint, startTime time.Time) error {
	b.store[key] = struct {
		count     uint
		startTime time.Time
	}{
		count:     count,
		startTime: startTime,
	}

	return nil
}

func (b *InMemoryucketList) Get(key string) (uint, time.Time, error) {
	if v, ok := b.store[key]; ok {
		return v.count, v.startTime, nil
	}

	return 0, time.Now(), nil
}

func (b *InMemoryucketList) Remove(key string) error {
	delete(b.store, key)

	return nil
}
