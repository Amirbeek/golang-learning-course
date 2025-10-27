package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound         = errors.New("record not found")
	NotRowEffectedError = errors.New("error rows affected")
	DeleteError         = errors.New("error deleting record")
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
	}
	Comments interface {
		GetByPostId(context.Context, int64) ([]*Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UserStore{db},
		Comments: &CommentsStore{db},
	}
}
