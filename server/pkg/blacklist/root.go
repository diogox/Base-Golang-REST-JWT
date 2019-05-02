package blacklist

import (
	"github.com/go-redis/redis"
	"time"
)

const failedLoginPrefix = "failed_login:blacklist:"

func NewBlacklist(host string, accountLockDuration int) (*Blacklist, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	return &Blacklist{
		client:              client,
		accountLockDuration: accountLockDuration,
	}, err
}

type Blacklist struct {
	client              *redis.Client
	accountLockDuration int
}

func (b *Blacklist) GetFailedLoginCountByUserID(userID string) (string, *time.Duration, error) {
	key := failedLoginPrefix + userID

	// Get count
	count, err := b.client.Get(key).Result()
	if err != nil {
		return "", nil, err
	}

	// get time left
	timeLeftUntilExpire, err := b.client.TTL(key).Result()
	return count, &timeLeftUntilExpire, err
}

func (b *Blacklist) IncrementFailedLoginCountByUserID(userID string) error {
	key := failedLoginPrefix + userID

	err := b.client.Incr(key).Err()
	if err != nil {
		return err
	}

	lockDuration := time.Duration(b.accountLockDuration) * time.Minute
	return b.client.Expire(key, lockDuration).Err()
}

func (b *Blacklist) ResetFailedLoginCountByUserID(userID string) error {
	key := failedLoginPrefix + userID
	return b.client.Del(key).Err()
}
