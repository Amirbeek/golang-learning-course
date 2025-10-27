package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserId    int64      `json:"userId"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Comments  []*Comment `json:"comments"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		p.Content,
		p.Title,
		p.UserId,
		pq.Array(p.Tags),
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (s *PostsStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `
SELECT id, content, title, user_id, tags, created_at, updated_at
FROM posts WHERE id = $1`
	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserId,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostsStore) DeleteById(ctx context.Context, id int64) error {
	query := `
    DELETE FROM posts WHERE id = $1;
`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return DeleteError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NotRowEffectedError
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) Update(ctx context.Context, p *Post) error {
	query := `
    UPDATE posts
    SET title = $1, content = $2
    WHERE id = $3
    `
	_, err := s.db.ExecContext(ctx, query, p.Title, p.Content, p.ID)

	return err
}
