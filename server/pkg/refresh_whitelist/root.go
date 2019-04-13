package refresh_whitelist

import (
	"github.com/go-redis/redis"
	"time"
)

func NewWhitelist(host string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	return client, err
}

func AddToWhitelist(client *redis.Client, tokenStr string, tokenDuration int) error {
	spareTime := 1
	expiresIn := time.Minute * time.Duration(tokenDuration + spareTime)
	err := client.Set(tokenStr, "", expiresIn).Err()
	if err != nil {
		return err
	}

	return nil
}