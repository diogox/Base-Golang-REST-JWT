package refresh_whitelist

import (
	"github.com/go-redis/redis"
	"time"
)

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

func (w *Whitelist) Set(key string, tokenStr string, tokenDuration int) error {
	spareTime := 1
	expiresIn := time.Minute * time.Duration(tokenDuration + spareTime)
	err := w.client.Set(key, tokenStr, expiresIn).Err()
	if err != nil {
		return err
	}

	return nil
}

func (w *Whitelist) Get(key string) (string, error) {
	return w.client.Get(key).Result()
}

func (w *Whitelist) Del(key string) (int64, error) {
	return w.client.Del(key).Result()
}