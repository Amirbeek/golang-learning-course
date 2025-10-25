package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"-"`
	CreatedAt time.Time `db:"created_at"`
}
type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, u *User) error {
	q := `
        INSERT INTO users (username, password, email)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	return s.db.QueryRowContext(ctx, q,
		u.Username,
		u.Password,
		u.Email,
	).Scan(&u.ID, &u.CreatedAt)
}
