package response

import (
	"time"

	"github.com/google/uuid"
)

type EditPostResponse struct {
	PostID             string `json:"postId"`
	PostCategoryID     string `json:"postCategoryId"`
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	Introduction       string `json:"introduction"`
	Content            string `json:"content"`
	MainImageAlt       string `json:"mainImageAlt"`
	MainImagePath      string `json:"mainImagePath"`
	Published          bool   `json:"published"`
	ThumbnailImagePath string `json:"thumbnailImagePath"`
	ThumbnailImageAlt  string `json:"thumbnailImageAlt"`
}

type UploadSunEditorResponse struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Size string `json:"size"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`	
}

type UserResponse struct {
	FirstName          string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
