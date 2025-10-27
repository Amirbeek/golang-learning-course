package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  password  `db:"-"`
	CreatedAt time.Time `db:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.MinCost)

	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	q := `
        INSERT INTO users (username, password, email)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, q, u.Username, u.Password.hash, u.Email).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		var pqe *pq.Error
		if errors.As(err, &pqe) && pqe.Code == "23505" {
			switch pqe.Constraint {
			case "users_email_key":
				return ErrDuplicateEmail
			case "users_username_key":
				return ErrDuplicateUsername
			}
		}
		return err
	}
	return nil
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

func (s *UserStore) CreateAndInvite(ctx context.Context, u *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.Create(ctx, tx, u); err != nil {
			return err
		}
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, u.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	const q = `INSERT INTO user_invitations (user_id, token, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	expiryTime := time.Now().Add(exp)

	_, err := tx.ExecContext(ctx, q, userID, token, expiryTime)
	return err
}
