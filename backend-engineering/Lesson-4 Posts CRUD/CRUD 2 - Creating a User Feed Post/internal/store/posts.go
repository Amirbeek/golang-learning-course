package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"userId"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
		&post.Tags,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
