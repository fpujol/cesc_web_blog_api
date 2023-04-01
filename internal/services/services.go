package services

import (
	"blogapi/api/request"
	db "blogapi/db/sqlc"
	"blogapi/pkg/password"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func InitUser(context context.Context, store db.Store) (db.User, error) {
	
	userFound, err := store.GetUserByEmail(context,"fpbosch@gmail.com")
	
	if err!=nil {
		hashedPassword, err := password.HashPassword("fpbosch@gmail.com","1234")
		args := db.CreateUserParams{
			FirstName: "Francesc",
			LastName: "Pujol",
			Email: "fpbosch@gmail.com",
			HashedPassword: hashedPassword,
			CreatedAt:          time.Now(),
			CreatedBy:          "Francesc Pujol",
		}
	
		result, err := store.CreateUser(context, args)
	
		return result, err
	}
	
	return userFound, err
}


func CreatePost(context context.Context, store db.Store, req request.RegisterNewPostRequest) (db.Post, error) {
	
	postCategoryId, errCategoryID := uuid.Parse(req.PostCategoryId)

	if errCategoryID!=nil {
		return db.Post{}, errors.New("no category")
	}

	args := db.CreatePostParams{
		PostID:             uuid.New(),
		Title:              req.Title,
		Slug:               slug.Make(req.Title),
		Introduction:       processNullOrString(req.Introduction),
		MainImageAlt:       processNullOrString(req.MainImageAlt),
		MainImagePath:      processNullOrString(req.MainImagePath),
		ThumbnailImageAlt:  processNullOrString(req.ThumbnailImageAlt),
		ThumbnailImagePath: processNullOrString(req.ThumbnailImagePath),
		PostCategoryID:     postCategoryId,
		Content:            processNullOrString(req.Content),
		Published:          req.Published,
		Author:             "Francesc Pujol",
		AuthorImagePath:    "../images/cesc.jpg",
		CreatedAt:          time.Now(),
		CreatedBy:          "Francesc Pujol",
	}

	result, err := store.CreatePost(context, args)

	return result, err
}

func UpdatePost(context context.Context, store db.Store, req request.EditPostRequest) (db.Post, error) {
	
	postCategoryId, errCategoryID := uuid.Parse(req.PostCategoryId)

	if errCategoryID!=nil {
		return db.Post{}, errors.New("no category")
	}

	postIdGuid, errPostId := uuid.Parse(req.PostId)

	if errPostId != nil {
		return db.Post{}, errors.New("no category")
	}

	args := db.UpdatePostParams{
		PostID:             postIdGuid,
		Title:              req.Title,
		Slug:               slug.Make(req.Title),
		Introduction:       processNullOrString(req.Introduction),
		MainImageAlt:       processNullOrString(req.MainImageAlt),
		ThumbnailImageAlt:  processNullOrString(req.ThumbnailImageAlt),
		PostCategoryID:     postCategoryId,
		Content:            processNullOrString(req.Content),
		Published:          req.Published,
		Author:             "Francesc Pujol",
		AuthorImagePath:    "../images/cesc.jpg",
		UpdatedAt:          time.Now(),
		UpdatedBy:          processNullOrString("Francesc Pujol"),
	}

	result, err := store.UpdatePost(context, args)

	return result, err
}

func UpdateMainImagePost(context context.Context, store db.Store, postId string, postMainImagePath string) (db.Post, error) {
	postIdGuid, errPostId := uuid.Parse(postId)

	if errPostId != nil {
		return db.Post{}, errors.New("no category")
	}

	args := db.UpdateMainImagePostParams{
		PostID:             postIdGuid,
		MainImagePath:      processNullOrString(postMainImagePath),
		UpdatedAt:          time.Now(),
		UpdatedBy:          processNullOrString("Francesc Pujol"),
	}

	result, err := store.UpdateMainImagePost(context, args)

	return result, err
}

func UpdateThubmnailImagePost(context context.Context, store db.Store, postId string, postThumbnailImagePath string) (db.Post, error) {
	postIdGuid, errPostId := uuid.Parse(postId)

	if errPostId != nil {
		return db.Post{}, errors.New("no category")
	}

	if errPostId != nil {
		return db.Post{}, errors.New("no category")
	}

	args := db.UpdateThumbnailImagePostParams{
		PostID:             postIdGuid,
		ThumbnailImagePath:      processNullOrString(postThumbnailImagePath),
		UpdatedAt:          time.Now(),
		UpdatedBy:          processNullOrString("Francesc Pujol"),
	}

	result, err := store.UpdateThumbnailImagePost(context, args)

	return result, err
}