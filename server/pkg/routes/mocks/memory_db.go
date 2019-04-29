package mocks

import "github.com/pkg/errors"

func NewMockInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		items: make(map[string]string, 0),
	}
}

type InMemoryDB struct {
	items map[string]string
}

func (im *InMemoryDB) Get(key string) (string, error) {
	for _key, _value := range im.items {
		if _key == key {
			return _value, nil
		}
	}

	return "", errors.New("item not found")
}

func (im *InMemoryDB) Set(key string, value string, valueDurationInMinutes int) error {
	im.items[key] = value
	return nil
}

func (im *InMemoryDB) Del(key string) (int64, error) {
	delete(im.items, key)
	if im.items[key] != "" {
		return 1, nil
	}

	return 0, nil
}