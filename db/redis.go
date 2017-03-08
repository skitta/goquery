package db

import (
	"gopkg.in/redis.v5"
)

// Redis pack a redis.Client
type Redis struct {
	Client *redis.Client
}

var redisClient *redis.Client

func init() {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
}

// Init redisClient
func (r *Redis) Init() {
	r.Client = redisClient
}

// Close redisClient
func (r *Redis) Close() error {
	return r.Client.Close()
}

// Index return all redis keys
func (r *Redis) Index() ([]string, error) {
	keys, err := r.Client.Keys("*").Result()
	return keys, err
}

// Update a record
func (r *Redis) Update(key, value string) error {
	err := r.Client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Del a record
func (r *Redis) Del(key string) error {
	err := r.Client.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}
