package blacklist

import (
	"github.com/go-redis/redis"
	"time"
)

const failedLoginPrefix = "failed_login:blacklist:"

func NewBlacklist(host string) (*Blacklist, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	return &Blacklist{
		client: client,
	}, err
}

type Blacklist struct {
	client *redis.Client
}

func (b *Blacklist) GetFailedLoginCountByUserID(userID string) (string, error) {
	return b.client.Get(failedLoginPrefix + userID).Result()
}

func (b *Blacklist) IncrementFailedLoginCountByUserID(userID string) error {
	err :=  b.client.Incr(failedLoginPrefix + userID).Err()
	if err != nil {
		return err
	}
	return b.client.Expire(failedLoginPrefix + userID, 10 * time.Minute).Err()
}

func (b *Blacklist) ResetFailedLoginCountByUserID(userID string) error {
	return b.client.Del(failedLoginPrefix + userID).Err()
}

