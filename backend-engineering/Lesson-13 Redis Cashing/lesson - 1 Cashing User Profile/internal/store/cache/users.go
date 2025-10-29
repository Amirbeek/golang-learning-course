package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/amirbeek/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err = json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, userID int64, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%d", userID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, string(json), UserExpTime).Err()
}
