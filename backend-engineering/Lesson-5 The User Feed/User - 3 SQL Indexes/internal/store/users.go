package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (s *UserStore) GetUserById(ctx context.Context, id int64) (*User, error) {
	q := `
	SELECT id, username, email, password, created_at 
	FROM users
	WHERE id = $1	
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var user User

	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) DeleteById(ctx context.Context, id int64) error {
	const q = `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete users id=%d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected users id=%d: %w", id, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
