// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetPost(ctx context.Context, slug string) (Post, error)
	GetPostComments(ctx context.Context, postID uuid.UUID) ([]Comment, error)
	GetPublicPosts(ctx context.Context, published bool) ([]Post, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListPosts(ctx context.Context) ([]Post, error)
	UpdateMainImagePost(ctx context.Context, arg UpdateMainImagePostParams) (Post, error)
	UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error)
	UpdateThumbnailImagePost(ctx context.Context, arg UpdateThumbnailImagePostParams) (Post, error)
	UpdateUserByEmail(ctx context.Context, arg UpdateUserByEmailParams) (User, error)
}

var _ Querier = (*Queries)(nil)
