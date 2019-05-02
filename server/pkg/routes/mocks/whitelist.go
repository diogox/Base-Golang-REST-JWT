package mocks

import "github.com/pkg/errors"

const refreshTokenPrefix = "refresh:"
const resetPasswordTokenPrefix = "reset:"

func NewWhitelist() *Whitelist {
	return &Whitelist{
		items: make(map[string]string, 0),
	}
}

type Whitelist struct {
	items map[string]string
}

// Refresh Token
func (im *Whitelist) GetRefreshTokenByUserID(key string) (string, error) {
	key = refreshTokenPrefix + key

	for _key, _value := range im.items {
		if _key == key {
			return _value, nil
		}
	}

	return "", errors.New("item not found")
}

func (im *Whitelist) SetRefreshTokenByUserID(key string, value string, valueDurationInMinutes int) error {
	key = refreshTokenPrefix + key

	im.items[key] = value
	return nil
}

func (im *Whitelist) DelRefreshTokenByUserID(key string) (int64, error) {
	key = refreshTokenPrefix + key

	delete(im.items, key)
	if im.items[key] != "" {
		return 1, nil
	}

	return 0, nil
}

// Reset Password Token
func (im *Whitelist) GetResetPasswordTokenByUserID(key string) (string, error) {
	key = resetPasswordTokenPrefix + key

	for _key, _value := range im.items {
		if _key == key {
			return _value, nil
		}
	}

	return "", errors.New("item not found")
}

func (im *Whitelist) SetResetPasswordTokenByUserID(key string, value string, valueDurationInMinutes int) error {
	key = resetPasswordTokenPrefix + key

	im.items[key] = value
	return nil
}

func (im *Whitelist) DelResetPasswordTokenByUserID(key string) (int64, error) {
	key = resetPasswordTokenPrefix + key

	delete(im.items, key)
	if im.items[key] != "" {
		return 1, nil
	}

	return 0, nil
}