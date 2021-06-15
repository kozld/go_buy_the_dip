package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var (
	key      = "assets"
	ctx      = context.Background()
	redisPwd = os.Getenv("REDIS_PASSWORD")
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() Store {
	client := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: redisPwd, // no password set
		DB:       1,        // use default DB
	})

	return &RedisStore{client}
}

func (rs *RedisStore) Set(assets []Asset) error {

	asBytes, _ := json.Marshal(assets)
	err := rs.client.Set(ctx, key, asBytes, 0).Err()
	if err != nil {
		return err
	}

	return err
}

func (rs *RedisStore) Get() ([]Asset, error) {

	val, err := rs.client.Get(ctx, key).Bytes()
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, nil
	}

	var assets []Asset
	err = json.Unmarshal(val, &assets)

	if len(assets) == 0 {
		return nil, nil
	}

	return assets, err
}
