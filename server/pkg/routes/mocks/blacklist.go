package mocks

import (
	"github.com/pkg/errors"
	"strconv"
	"time"
)

const failedLoginPrefix = "failed_login:"

func NewBlacklist() *Blacklist {
	return &Blacklist{
		items: make(map[string]string, 0),
	}
}

type Blacklist struct {
	items map[string]string
}

// Failed Login Blacklist
func (im *Blacklist) GetFailedLoginCountByUserID(userID string) (string, *time.Duration, error) {
	key := failedLoginPrefix + userID

	for _key, _value := range im.items {
		if _key == key {
			return _value, nil, nil
		}
	}

	return "", nil, errors.New("item not found")
}

func (im *Blacklist) IncrementFailedLoginCountByUserID(userID string) error {
	key := refreshTokenPrefix + userID

	previousCount := im.items[key]
	if previousCount == "" {
		previousCount = "0"
	}

	previousCountInt, _ := strconv.Atoi(previousCount)
	im.items[key] = strconv.Itoa(previousCountInt+1)

	return nil
}

func (im *Blacklist) ResetFailedLoginCountByUserID(userID string) error {
	key := refreshTokenPrefix + userID
	im.items[key] = "0"

	return nil
}