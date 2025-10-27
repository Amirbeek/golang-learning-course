package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record not found")
	NotRowEffectedError  = errors.New("error rows affected")
	DeleteError          = errors.New("error deleting record")
	ErrEditConflict      = errors.New("edit conflict")
	QueryTimeOutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		DeleteById(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetUserById(context.Context, int64) (*User, error)
		DeleteById(context.Context, int64) error
	}
	Comments interface {
		GetByPostId(context.Context, int64) ([]*Comment, error)
		Create(context.Context, *Comment) error
	}
	Follower interface {
		Follow(ctx context.Context, followerId, userId int64) error
		UnFollow(ctx context.Context, followerId, userId int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UserStore{db},
		Comments: &CommentsStore{db},
		Follower: &FollowerStore{db},
	}
}
