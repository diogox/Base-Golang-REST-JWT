package mocks

import "github.com/pkg/errors"

func NewMockInMemoryDB() *InMemoryDB {
	return &InMemoryDB{}
}

type InMemoryDB struct {
	items []string
}

func (im *InMemoryDB) Get(tokenStr string) (string, error) {
	for _, item := range im.items {
		if item == tokenStr {
			return item, nil
		}
	}

	return "", errors.New("item not found")
}

func (im *InMemoryDB) Set(value string, valueDurationInMinutes int) error {
	for _, item := range im.items {
		if item == value {
			return errors.New("already exists")
		}
	}

	im.items = append(im.items, value)
	return nil
}

func (im *InMemoryDB) Del(tokenStr string) (int64, error) {
	for i, item := range im.items {
		if item == tokenStr {
			im.items = append(im.items[:i], im.items[i+1:]...)
			return 1, nil
		}
	}

	return 0, errors.New("not found")
}