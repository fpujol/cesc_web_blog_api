package api

import (
	"blogapi/api/response"
	db "blogapi/db/sqlc"
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func Test_MapPost(t *testing.T) {
	post := db.Post{
		PostID: uuid.New(),
		Title:              "Hello World",
		Slug: "hello-world",
		Introduction:       sql.NullString{String: "Me and you!", Valid: true},
		MainImageAlt:       sql.NullString{String: "Alt 1", Valid: true},
		MainImagePath:      sql.NullString{String: "img1", Valid: true},
		ThumbnailImageAlt:  sql.NullString{String: "Alt 2", Valid: true},
		ThumbnailImagePath: sql.NullString{String: "img2", Valid: true},
		Content:            sql.NullString{String: "<p>No way</p>", Valid: true},
		Published:          false,
		Author:             "Francesc Pujol",
		AuthorImagePath:    "../images/cesc.jpg",
	}

	type args struct {
		post db.Post
	}

	testTable := []struct {
		name string
		args args
		want response.EditPostResponse
	}{
		{
			name:"DeepEqual",
			args: args{
				post: post,
			},
			want: mapPostToResponse(post),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapPostToResponse(tt.args.post); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapPost() = %v, want %v", got, tt.want)
			}
		})
	}
}
