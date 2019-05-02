package whitelist

import (
	"github.com/go-redis/redis"
	"time"
)

const refreshTokenPrefix = "whitelist:refresh:"
const resetPasswordTokenPrefix = "whitelist:reset:"

func NewWhitelist(host string) (*Whitelist, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	return &Whitelist{
		client: client,
	}, err
}

type Whitelist struct {
	client *redis.Client
}

// Refresh Token
func (w *Whitelist) SetRefreshTokenByUserID(key string, tokenStr string, tokenDuration int) error {
	spareTime := 1
	expiresIn := time.Minute * time.Duration(tokenDuration + spareTime)
	err := w.client.Set(refreshTokenPrefix + key, tokenStr, expiresIn).Err()
	if err != nil {
		return err
	}

	return nil
}

func (w *Whitelist) GetRefreshTokenByUserID(key string) (string, error) {
	return w.client.Get(refreshTokenPrefix + key).Result()
}

func (w *Whitelist) DelRefreshTokenByUserID(key string) (int64, error) {
	return w.client.Del(refreshTokenPrefix + key).Result()
}

// Reset Password Token
func (w *Whitelist) SetResetPasswordTokenByUserID(key string, tokenStr string, tokenDuration int) error {
	spareTime := 1
	expiresIn := time.Minute * time.Duration(tokenDuration + spareTime)
	err := w.client.Set(resetPasswordTokenPrefix + key, tokenStr, expiresIn).Err()
	if err != nil {
		return err
	}

	return nil
}

func (w *Whitelist) GetResetPasswordTokenByUserID(key string) (string, error) {
	return w.client.Get(resetPasswordTokenPrefix + key).Result()
}

func (w *Whitelist) DelResetPasswordTokenByUserID(key string) (int64, error) {
	return w.client.Del(resetPasswordTokenPrefix + key).Result()
}