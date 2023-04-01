package api

import (
	"blogapi/api/response"
	db "blogapi/db/sqlc"
)

func mapPostToResponse(post db.Post) response.EditPostResponse {
	return response.EditPostResponse{
		PostID:       post.PostID.String(),
		Title:        post.Title,
		Slug:         post.Slug,
		PostCategoryID: post.PostCategoryID.String(),
		Introduction: post.Introduction.String,
		MainImageAlt:       post.MainImageAlt.String,
		MainImagePath:      post.MainImagePath.String, 
		ThumbnailImageAlt:  post.ThumbnailImageAlt.String,
		ThumbnailImagePath: post.ThumbnailImagePath.String,
		Content:      post.Content.String,
		Published: post.Published,
	}
}

func mapPostsToResponse(posts []db.Post) []response.EditPostResponse {

	if len(posts)>0 {
		p:=make([]response.EditPostResponse,0)
		for _, v := range posts {
			p = append(p,mapPostToResponse(v))
		}
		return p;
	}
	return []response.EditPostResponse{}
}


func mapUserToResponse(user db.User) response.UserResponse {
    return response.UserResponse{
        FirstName:          user.FirstName,
        LastName:          user.LastName,
        Email:             user.Email,
        PasswordChangedAt: user.PasswordChangedAt,
        CreatedAt:         user.CreatedAt,
    }
}