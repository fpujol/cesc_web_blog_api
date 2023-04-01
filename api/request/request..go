package request

type RegisterNewPostRequest struct {
	PostCategoryId     string `json:"postCategoryId" binding:"required"`
	Title              string `json:"title" binding:"required"`
	Introduction       string `json:"introduction"`
	Content            string `json:"content"`
	MainImageAlt       string `json:"mainImageAlt"`
	MainImagePath      string `json:"mainImagePath"`
	ThumbnailImagePath string `json:"thumbnailImagePath"`
	ThumbnailImageAlt  string `json:"thumbnailImageAlt"`
	Published          bool   `json:"published"`
}

type EditPostRequest struct {
	PostId            string `json:"postId" binding:"required"`
	PostCategoryId    string `json:"postCategoryId" binding:"required"`
	Title             string `json:"title" binding:"required"`
	Introduction      string `json:"introduction"`
	Content           string `json:"content"`
	MainImageAlt      string `json:"mainImageAlt"`
	ThumbnailImageAlt string `json:"thumbnailImageAlt"`
	Published         bool   `json:"published"`
}

type EditMainImageRequest struct {
	PostId        string `json:"postId" binding:"required"`
	MainImagePath string `json:"mainImagePath"`
}

type EditThumbnailImageRequest struct {
	PostId             string `json:"postId" binding:"required"`
	ThumbnailImagePath string `json:"thumbnailImagePath"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=4"`
}