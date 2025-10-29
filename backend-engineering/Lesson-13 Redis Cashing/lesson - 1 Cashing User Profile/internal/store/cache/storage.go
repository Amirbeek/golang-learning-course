package cache

import (
	"context"

	"github.com/amirbeek/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Users interface {
		Get(ctx context.Context, userID int64) (*store.User, error)
		Set(ctx context.Context, userID int64, user *store.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
