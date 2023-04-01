package db

import (
	"blogapi/pkg/utils"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/require"
)

func createRandomPost(t *testing.T) Post {
	postId := uuid.New()
	categoryId, _ := uuid.Parse("4e86fd8f-c75b-498d-982d-48f7033a3a47")
	title:= utils.RandomString(6)

	arg := CreatePostParams{
		PostID:             postId,
		Title:              title,
		Slug:               slug.Make(title),
		Introduction:       sql.NullString{String: utils.RandomString(6), Valid: true},
		MainImageAlt:       sql.NullString{String: utils.RandomString(6), Valid: true},
		MainImagePath:      sql.NullString{String: utils.RandomString(6), Valid: true},
		ThumbnailImageAlt:  sql.NullString{String: utils.RandomString(6), Valid: true},
		ThumbnailImagePath: sql.NullString{String: utils.RandomString(6), Valid: true},
		Content:            sql.NullString{String: utils.RandomString(6), Valid: true},
		PostCategoryID:     categoryId,
		Author:             "Francesc Pujol",
		AuthorImagePath:    "../images/cesc.jpg",
		CreatedAt:          time.Now(),
		CreatedBy:          "Francesc Pujol",
	}
	
	post, err := testQueries.CreatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	return post;
}

func createRandomPublishedPost(t *testing.T) Post {
	postId := uuid.New()
	categoryId, _ := uuid.Parse("4e86fd8f-c75b-498d-982d-48f7033a3a47")
	title:= utils.RandomString(6)

	arg := CreatePostParams{
		PostID:             postId,
		Title:              title,
		Slug:               slug.Make(title),
		Introduction:       sql.NullString{String: utils.RandomString(6), Valid: true},
		MainImageAlt:       sql.NullString{String: utils.RandomString(6), Valid: true},
		MainImagePath:      sql.NullString{String: utils.RandomString(6), Valid: true},
		ThumbnailImageAlt:  sql.NullString{String: utils.RandomString(6), Valid: true},
		ThumbnailImagePath: sql.NullString{String: utils.RandomString(6), Valid: true},
		Content:            sql.NullString{String: utils.RandomString(6), Valid: true},
		PostCategoryID:     categoryId,
		Published:			true,
		Author:             "Francesc Pujol",
		AuthorImagePath:    "../images/cesc.jpg",
		CreatedAt:          time.Now(),
		CreatedBy:          "Francesc Pujol",
	}
	
	post, err := testQueries.CreatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	return post;
}

func TestCreatePost(t *testing.T) {
	createRandomPost(t)
}

func TestGetPost(t *testing.T) {
	post1 := createRandomPost(t)
	post2, err := testQueries.GetPost(context.Background(), post1.Slug)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.PostID, post2.PostID)
	require.Equal(t, post1.Title, post2.Title)
	require.Equal(t, post1.Slug, post2.Slug)
	require.WithinDuration(t, post1.CreatedAt, post2.CreatedAt, time.Second)	
}

func TestGetPublicPosts(t *testing.T) {
	var lastPost Post
	for i := 0; i < 10; i++ {
		lastPost = createRandomPublishedPost(t)
	}

	posts, err := testQueries.GetPublicPosts(context.Background(), true)
	require.NoError(t, err)
	require.NotEmpty(t, posts)

	for _, post := range posts {
		require.NotEmpty(t, post)
		require.Equal(t, lastPost.Author, post.Author)
	}
}

func TestListPosts(t *testing.T) {
	var lastPost Post
	for i := 0; i < 10; i++ {
		lastPost = createRandomPost(t)
	}

	posts, err := testQueries.ListPosts(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, posts)

	for _, post := range posts {
		require.NotEmpty(t, post)
		require.Equal(t, lastPost.Author, post.Author)
	}
}

func TestUpdatePost(t *testing.T) {
	title:= utils.RandomString(6)
	post1 := createRandomPost(t)
	categoryId, _ := uuid.Parse("4e86fd8f-c75b-498d-982d-48f7033a3a47")
		
	arg := UpdatePostParams{
		PostID:      post1.PostID,
		Title: title,
		Slug: slug.Make(title),
		PostCategoryID: categoryId,
	}

	post2, err := testQueries.UpdatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.PostID, post2.PostID)
	require.Equal(t, title, post2.Title)
	require.Equal(t, slug.Make(title), post2.Slug)
	require.WithinDuration(t, post1.CreatedAt, post2.CreatedAt, time.Second)
}
