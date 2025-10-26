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
	Version   int        `json:"version"`
	User      *User      `json:"user"`
}

type PostsStore struct {
	db *sql.DB
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

func (s *PostsStore) GetUserFeed(ctx context.Context, id int64) ([]PostWithMetadata, error) {
	query := `
		SELECT 
			p.id,
			p.user_id,
			p.title,
			p.content,
			p.created_at,
			p.version,
			p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
		OR p.user_id IN (
			SELECT user_id FROM followers WHERE follower_id = $1
		)
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		p.User = &User{} // ✅ fix nil pointer

		err := rows.Scan(
			&p.ID,
			&p.UserId,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feed, nil
}

func (s *PostsStore) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

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
		SELECT id, content, title, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserId,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

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
		SET title = $1, content = $2, version = version + 1, updated_at = NOW()
		WHERE id = $3 AND version = $4
		RETURNING version;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, p.Title, p.Content, p.ID, p.Version).Scan(&p.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
