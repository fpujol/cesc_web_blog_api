// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: posts.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (post_id,
                  slug,
                  title,
                  introduction,
                  post_category_id,
                  main_image_alt, 
                  main_image_path, 
                  thumbnail_image_alt, 
                  thumbnail_image_path, 
                  content, 
                  author, 
                  author_image_path, 
                  author_image_alt,
                  published,
                  published_at,
                  published_by, 
                  created_at, 
                  created_by) 
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) 
RETURNING post_id, post_category_id, title, slug, introduction, content, main_image_alt, main_image_path, thumbnail_image_path, thumbnail_image_alt, author, author_image_path, author_image_alt, published, published_at, published_by, created_at, created_by, updated_at, updated_by
`

type CreatePostParams struct {
	PostID             uuid.UUID      `json:"post_id"`
	Slug               string         `json:"slug"`
	Title              string         `json:"title"`
	Introduction       sql.NullString `json:"introduction"`
	PostCategoryID     uuid.UUID      `json:"post_category_id"`
	MainImageAlt       sql.NullString `json:"main_image_alt"`
	MainImagePath      sql.NullString `json:"main_image_path"`
	ThumbnailImageAlt  sql.NullString `json:"thumbnail_image_alt"`
	ThumbnailImagePath sql.NullString `json:"thumbnail_image_path"`
	Content            sql.NullString `json:"content"`
	Author             string         `json:"author"`
	AuthorImagePath    string         `json:"author_image_path"`
	AuthorImageAlt     string         `json:"author_image_alt"`
	Published          bool           `json:"published"`
	PublishedAt        sql.NullTime   `json:"published_at"`
	PublishedBy        sql.NullString `json:"published_by"`
	CreatedAt          time.Time      `json:"created_at"`
	CreatedBy          string         `json:"created_by"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.PostID,
		arg.Slug,
		arg.Title,
		arg.Introduction,
		arg.PostCategoryID,
		arg.MainImageAlt,
		arg.MainImagePath,
		arg.ThumbnailImageAlt,
		arg.ThumbnailImagePath,
		arg.Content,
		arg.Author,
		arg.AuthorImagePath,
		arg.AuthorImageAlt,
		arg.Published,
		arg.PublishedAt,
		arg.PublishedBy,
		arg.CreatedAt,
		arg.CreatedBy,
	)
	var i Post
	err := row.Scan(
		&i.PostID,
		&i.PostCategoryID,
		&i.Title,
		&i.Slug,
		&i.Introduction,
		&i.Content,
		&i.MainImageAlt,
		&i.MainImagePath,
		&i.ThumbnailImagePath,
		&i.ThumbnailImageAlt,
		&i.Author,
		&i.AuthorImagePath,
		&i.AuthorImageAlt,
		&i.Published,
		&i.PublishedAt,
		&i.PublishedBy,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const getPost = `-- name: GetPost :one
SELECT post_id, post_category_id, title, slug, introduction, content, main_image_alt, main_image_path, thumbnail_image_path, thumbnail_image_alt, author, author_image_path, author_image_alt, published, published_at, published_by, created_at, created_by, updated_at, updated_by FROM posts 
WHERE slug=$1 LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, slug string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, slug)
	var i Post
	err := row.Scan(
		&i.PostID,
		&i.PostCategoryID,
		&i.Title,
		&i.Slug,
		&i.Introduction,
		&i.Content,
		&i.MainImageAlt,
		&i.MainImagePath,
		&i.ThumbnailImagePath,
		&i.ThumbnailImageAlt,
		&i.Author,
		&i.AuthorImagePath,
		&i.AuthorImageAlt,
		&i.Published,
		&i.PublishedAt,
		&i.PublishedBy,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const getPublicPosts = `-- name: GetPublicPosts :many
SELECT post_id, post_category_id, title, slug, introduction, content, main_image_alt, main_image_path, thumbnail_image_path, thumbnail_image_alt, author, author_image_path, author_image_alt, published, published_at, published_by, created_at, created_by, updated_at, updated_by FROM posts
WHERE published=$1 order by published_at
`

func (q *Queries) GetPublicPosts(ctx context.Context, published bool) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPublicPosts, published)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.PostID,
			&i.PostCategoryID,
			&i.Title,
			&i.Slug,
			&i.Introduction,
			&i.Content,
			&i.MainImageAlt,
			&i.MainImagePath,
			&i.ThumbnailImagePath,
			&i.ThumbnailImageAlt,
			&i.Author,
			&i.AuthorImagePath,
			&i.AuthorImageAlt,
			&i.Published,
			&i.PublishedAt,
			&i.PublishedBy,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPosts = `-- name: ListPosts :many
SELECT post_id, post_category_id, title, slug, introduction, content, main_image_alt, main_image_path, thumbnail_image_path, thumbnail_image_alt, author, author_image_path, author_image_alt, published, published_at, published_by, created_at, created_by, updated_at, updated_by FROM posts 
ORDER BY title
`

func (q *Queries) ListPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, listPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.PostID,
			&i.PostCategoryID,
			&i.Title,
			&i.Slug,
			&i.Introduction,
			&i.Content,
			&i.MainImageAlt,
			&i.MainImagePath,
			&i.ThumbnailImagePath,
			&i.ThumbnailImageAlt,
			&i.Author,
			&i.AuthorImagePath,
			&i.AuthorImageAlt,
			&i.Published,
			&i.PublishedAt,
			&i.PublishedBy,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMainImagePost = `-- name: UpdateMainImagePost :one
UPDATE posts
  set main_image_path=$2, 
  updated_at=$3,
  updated_by=$4
WHERE post_id = $1
RETURNING post_id, post_category_id, title, slug, introduction, content, main_image_alt, main_image_path, thumbnail_image_path, thumbnail_image_alt, author, author_image_path, author_image_alt, published, published_at, published_by, created_at, created_by, updated_at, updated_by
`

type UpdateMainImagePostParams struct {
	PostID        uuid.UUID      `json:"post_id"`
	MainImagePath sql.NullString `json:"main_image_path"`
	UpdatedAt     time.Time      `json:"updated_at"`
	UpdatedBy     sql.NullString `json:"updated_by"`
}

func (q *Queries) UpdateMainImagePost(ctx context.Context, arg UpdateMainImagePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updateMainImagePost,
		arg.PostID,
		arg.MainImagePath,
		arg.UpdatedAt,
		arg.UpdatedBy,
	)
	var i Post
	err := row.Scan(
		&i.PostID,
		&i.PostCategoryID,
		&i.Title,
		&i.Slug,
		&i.Introduction,
		&i.Content,
		&i.MainImageAlt,
		&i.MainImagePath,
		&i.ThumbnailImagePath,
		&i.ThumbnailImageAlt,
		&i.Author,
		&i.AuthorImagePath,
		&i.AuthorImageAlt,
		&i.Published,
		&i.PublishedAt,
		&i.PublishedBy,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const updatePost = `-- name: UpdatePost :one
UPDATE posts
  set post_category_id=$2,
  slug=$3,
  title=$4, 
  introduction=$5,
  main_image_alt=$6, 
  thumbnail_image_alt=$7, 
  content=$8,
  author=$9, 
  author_image_path=$10, 
  author_image_alt=$11,
  published=$12,
  published_at=$13,
  published_by=$14,
  updated_at=$15,
  updated_by=$16
WHERE post_id = $1
RETURNING post_id, post_category_id, title, slug, introduction, content, main_image_alt, main_image_path, thumbnail_image_path, thumbnail_image_alt, author, author_image_path, author_image_alt, published, published_at, published_by, created_at, created_by, updated_at, updated_by
`

type UpdatePostParams struct {
	PostID            uuid.UUID      `json:"post_id"`
	PostCategoryID    uuid.UUID      `json:"post_category_id"`
	Slug              string         `json:"slug"`
	Title             string         `json:"title"`
	Introduction      sql.NullString `json:"introduction"`
	MainImageAlt      sql.NullString `json:"main_image_alt"`
	ThumbnailImageAlt sql.NullString `json:"thumbnail_image_alt"`
	Content           sql.NullString `json:"content"`
	Author            string         `json:"author"`
	AuthorImagePath   string         `json:"author_image_path"`
	AuthorImageAlt    string         `json:"author_image_alt"`
	Published         bool           `json:"published"`
	PublishedAt       sql.NullTime   `json:"published_at"`
	PublishedBy       sql.NullString `json:"published_by"`
	UpdatedAt         time.Time      `json:"updated_at"`
	UpdatedBy         sql.NullString `json:"updated_by"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updatePost,
		arg.PostID,
		arg.PostCategoryID,
		arg.Slug,
		arg.Title,
		arg.Introduction,
		arg.MainImageAlt,
		arg.ThumbnailImageAlt,
		arg.Content,
		arg.Author,
		arg.AuthorImagePath,
		arg.AuthorImageAlt,
		arg.Published,
		arg.PublishedAt,
		arg.PublishedBy,
		arg.UpdatedAt,
		arg.UpdatedBy,
	)
	var i Post
	err := row.Scan(
		&i.PostID,
		&i.PostCategoryID,
		&i.Title,
		&i.Slug,
		&i.Introduction,
		&i.Content,
		&i.MainImageAlt,
		&i.MainImagePath,
		&i.ThumbnailImagePath,
		&i.ThumbnailImageAlt,
		&i.Author,
		&i.AuthorImagePath,
		&i.AuthorImageAlt,
		&i.Published,
		&i.PublishedAt,
		&i.PublishedBy,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const updateThumbnailImagePost = `-- name: UpdateThumbnailImagePost :one
UPDATE posts
  set thumbnail_image_path=$2, 
  updated_at=$3,
  updated_by=$4
WHERE post_id = $1
RETURNING post_id, post_category_id, title, slug, introduction, content, main_image_alt, main_image_path, thumbnail_image_path, thumbnail_image_alt, author, author_image_path, author_image_alt, published, published_at, published_by, created_at, created_by, updated_at, updated_by
`

type UpdateThumbnailImagePostParams struct {
	PostID             uuid.UUID      `json:"post_id"`
	ThumbnailImagePath sql.NullString `json:"thumbnail_image_path"`
	UpdatedAt          time.Time      `json:"updated_at"`
	UpdatedBy          sql.NullString `json:"updated_by"`
}

func (q *Queries) UpdateThumbnailImagePost(ctx context.Context, arg UpdateThumbnailImagePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updateThumbnailImagePost,
		arg.PostID,
		arg.ThumbnailImagePath,
		arg.UpdatedAt,
		arg.UpdatedBy,
	)
	var i Post
	err := row.Scan(
		&i.PostID,
		&i.PostCategoryID,
		&i.Title,
		&i.Slug,
		&i.Introduction,
		&i.Content,
		&i.MainImageAlt,
		&i.MainImagePath,
		&i.ThumbnailImagePath,
		&i.ThumbnailImageAlt,
		&i.Author,
		&i.AuthorImagePath,
		&i.AuthorImageAlt,
		&i.Published,
		&i.PublishedAt,
		&i.PublishedBy,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}
