package persistence

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

// RedisAdapter a simple adapter
type RedisAdapter struct {
	client *redis.Client
}

// Delete delete a key
func (r RedisAdapter) Delete(counterKey string) error {
	_, err := r.client.Del(counterKey).Result()
	return err
}

// IncrementCounter increment a counter
func (r RedisAdapter) IncrementCounter(counterKey string) (uint64, error) {

	result := r.client.Incr(fmt.Sprintf("%s:count", counterKey))

	if result.Err() != nil {
		return 0, result.Err()
	}

	return uint64(result.Val()), nil
}

// ReadCounter read a counter
func (r RedisAdapter) ReadCounter(counterKey string) (uint64, error) {
	response := r.client.Get(counterKey)

	if response == nil {
		return 0, errors.New("Persistence read problem")
	}

	if response.Err() == redis.Nil {
		return 0, errors.New("Counter not found")
	}

	return response.Uint64()
}

// Lookup perform a lookup on the redis db
func (r RedisAdapter) Lookup(key string) (string, error) {
	response, err := r.client.Get(key).Result()

	if err == redis.Nil {
		return "", errors.New("Key not found")
	} else if err != nil {
		return "", err
	} else {
		return response, nil
	}
}

// Persist persist a key value pair
func (r RedisAdapter) Persist(key string, value interface{}) error {
	err := r.client.Set(key, value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

// CreateRedisAdapter initalize new redis adapter
func CreateRedisAdapter() *RedisAdapter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisAdapter{
		client: rdb,
	}
}
