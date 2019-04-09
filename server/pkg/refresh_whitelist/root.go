package refresh_whitelist

import "github.com/go-redis/redis"

func NewWhitelist() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	return client, err
}
