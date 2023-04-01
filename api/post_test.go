package api

import (
	"blogapi/api/response"
	mockdb "blogapi/db/mock"
	db "blogapi/db/sqlc"
	"blogapi/pkg/token"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreatePostParamsMatcher struct {
	arg       db.CreatePostParams
	PostID    uuid.UUID
	CreatedAt time.Time
}

func (e eqCreatePostParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreatePostParams)
	if !ok {
		return false
	}

	if e.arg.Title != arg.Title ||
		e.arg.Introduction != arg.Introduction ||
		e.arg.Content != arg.Content {
		return false
	}

	e.arg.PostID = arg.PostID
	e.arg.CreatedAt = arg.CreatedAt
	result := reflect.DeepEqual(e.arg, arg)
	return result
}

func (e eqCreatePostParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and PostId %v and CreatedAt %v ", e.arg, e.PostID, e.CreatedAt)
}

func EqCreatePostParams(arg db.CreatePostParams, postId uuid.UUID, createdAt time.Time) gomock.Matcher {
	return eqCreatePostParamsMatcher{arg, postId, createdAt}
}

type eqEditPostParamsMatcher struct {
	arg       db.UpdatePostParams
	PostID    uuid.UUID
	UpdatedAt time.Time
}

func (e eqEditPostParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdatePostParams)
	if !ok {
		return false
	}

	if e.arg.Title != arg.Title ||
		e.arg.Introduction != arg.Introduction ||
		e.arg.Content != arg.Content {
		return false
	}

	e.arg.PostID = arg.PostID
	e.arg.UpdatedAt = arg.UpdatedAt
	result := reflect.DeepEqual(e.arg, arg)
	return result
}

func (e eqEditPostParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and PostId %v and UpdatedAt %v ", e.arg, e.PostID, e.UpdatedAt)
}

func EqEditPostParams(arg db.UpdatePostParams, postId uuid.UUID, createdAt time.Time) gomock.Matcher {
	return eqEditPostParamsMatcher{arg, postId, createdAt}
}

func requireBodyMatchPost(t *testing.T, body *bytes.Buffer, post response.EditPostResponse) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotPost response.EditPostResponse
	err = json.Unmarshal(data, &gotPost)
	require.NoError(t, err)
	require.Equal(t, post, gotPost)
}

func TestGetPostAPI(t *testing.T) {
	post := db.Post{
		PostID: uuid.New(),
		Title:              "Hello World",
		Slug:               "hello-world",
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

	postResponse := mapPostToResponse(post)
	
	testTable := []struct {
		name string
		slug string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		want func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			slug: "hello-world",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Eq(post.Slug)).
					Times(1).
					Return(post, nil)
			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "NotFound",
			slug: post.Slug,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Eq(post.Slug)).
					Times(1).
					Return(db.Post{}, sql.ErrNoRows)
			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			slug: post.Slug,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Eq(post.Slug)).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)
			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidSlug",
			slug: "",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testTable {
		tt := testTable[i]

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/post/%s", tt.slug)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tt.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tt.want(t, recorder)
		})
	}
}

func TestRegisterNewPost(t *testing.T) {
	postId := uuid.New()
	createdAt := time.Now()
	categoryId, _ := uuid.Parse("4e86fd8f-c75b-498d-982d-48f7033a3a47")

	post := db.Post{
		PostID: postId,
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

	postResponse := mapPostToResponse(post)
	
	testTable := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		want func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			body: gin.H{
				"title":              "Hello World",
				"slug": "hello-world",
				"introduction":       "Me and you!",
				"mainImageAlt":       "Alt 1",
				"mainImagePath":      "img1",
				"thumbnailImageAlt":  "Alt 2",
				"thumbnailImagePath": "img2",
				"content":            "<p>No way</p>",
				"postCategoryID":     categoryId.String(),
				"published":          false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreatePostParams{
					PostID:             postId,
					Title:              "Hello World",
					Slug:               "hello-world",
					Introduction:       sql.NullString{String: "Me and you!", Valid: true},
					MainImageAlt:       sql.NullString{String: "Alt 1", Valid: true},
					MainImagePath:      sql.NullString{String: "img1", Valid: true},
					ThumbnailImageAlt:  sql.NullString{String: "Alt 2", Valid: true},
					ThumbnailImagePath: sql.NullString{String: "img2", Valid: true},
					Content:            sql.NullString{String: "<p>No way</p>", Valid: true},
					PostCategoryID:     categoryId,
					Author:             "Francesc Pujol",
					AuthorImagePath:    "../images/cesc.jpg",
					CreatedAt:          createdAt,
					CreatedBy:          "Francesc Pujol",
				}

				store.EXPECT().
					CreatePost(gomock.Any(), EqCreatePostParams(arg, postId, createdAt)).
					Times(1).
					Return(post, nil)

			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "Required Title",
			body: gin.H{
				"introduction":       "Me and you!",
				"mainImageAlt":       "Alt 1",
				"mainImagePath":      "img1",
				"thumbnailImageAlt":  "Alt 2",
				"thumbnailImagePath": "img2",
				"content":            "<p>No way</p>",
				"postCategoryID":     categoryId.String(),
				"published":          false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreatePostParams{
					PostID:             postId,
					Title:              "Hello World",
					Slug:               "hello-world",
					Introduction:       sql.NullString{String: "Me and you!", Valid: true},
					MainImageAlt:       sql.NullString{String: "Alt 1", Valid: true},
					MainImagePath:      sql.NullString{String: "img1", Valid: true},
					ThumbnailImageAlt:  sql.NullString{String: "Alt 2", Valid: true},
					ThumbnailImagePath: sql.NullString{String: "img2", Valid: true},
					Content:            sql.NullString{String: "<p>No way</p>", Valid: true},
					PostCategoryID:     categoryId,
					Author:             "Francesc Pujol",
					AuthorImagePath:    "../images/cesc.jpg",
					CreatedAt:          createdAt,
					CreatedBy:          "Francesc Pujol",
				}

				store.EXPECT().
					CreatePost(gomock.Any(), EqCreatePostParams(arg, postId, createdAt)).
					Times(0)

			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Duplicate Title",
			body: gin.H{
				"title":              "Hello World",
				"introduction":       "Me and you!",
				"mainImageAlt":       "Alt 1",
				"mainImagePath":      "img1",
				"thumbnailImageAlt":  "Alt 2",
				"thumbnailImagePath": "img2",
				"content":            "<p>No way</p>",
				"postCategoryID":     categoryId.String(),
				"published":          false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreatePostParams{
					PostID:             postId,
					Title:              "Hello World",
					Slug:               "hello-world",
					Introduction:       sql.NullString{String: "Me and you!", Valid: true},
					MainImageAlt:       sql.NullString{String: "Alt 1", Valid: true},
					MainImagePath:      sql.NullString{String: "img1", Valid: true},
					ThumbnailImageAlt:  sql.NullString{String: "Alt 2", Valid: true},
					ThumbnailImagePath: sql.NullString{String: "img2", Valid: true},
					Content:            sql.NullString{String: "<p>No way</p>", Valid: true},
					PostCategoryID:     categoryId,
					Author:             "Francesc Pujol",
					AuthorImagePath:    "../images/cesc.jpg",
					CreatedAt:          createdAt,
					CreatedBy:          "Francesc Pujol",
				}

				store.EXPECT().
					CreatePost(gomock.Any(), EqCreatePostParams(arg, postId, createdAt)).
					Times(1).
					Return(db.Post{}, &pq.Error{Code: "23505"})

			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testTable {
		tt := testTable[i]

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			bodyData, err := json.Marshal(tt.body)
			require.NoError(t, err)

			url := "/api/v1/post"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyData))
			require.NoError(t, err)

			tt.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tt.want(t, recorder)
		})
	}
}

func TestEditPost(t *testing.T) {
	postId := uuid.New()
    updatedAt := time.Now()
	categoryId, _ := uuid.Parse("4e86fd8f-c75b-498d-982d-48f7033a3a47")

	post := db.Post{
		PostID: postId,
		Title:              "Hello World",
		Slug:               "hello-world",
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

	postResponse := mapPostToResponse(post)
	
	testTable := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		want func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			body: gin.H{
				"postId": postId,
				"title":              "Hello World",
				"slug":                 "hello-world",
				"introduction":       "Me and you!",
				"mainImageAlt":       "Alt 1",
				"mainImagePath":      "img1",
				"thumbnailImageAlt":  "Alt 2",
				"thumbnailImagePath": "img2",
				"content":            "<p>No way</p>",
				"postCategoryID":     categoryId.String(),
				"published":          false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostParams{
					PostID:             postId,
					Title:              "Hello World",
					Slug:               "hello-world",
					Introduction:       sql.NullString{String: "Me and you!", Valid: true},
					MainImageAlt:       sql.NullString{String: "Alt 1", Valid: true},
					ThumbnailImageAlt:  sql.NullString{String: "Alt 2", Valid: true},
					Content:            sql.NullString{String: "<p>No way</p>", Valid: true},
					PostCategoryID:     categoryId,
					Author:             "Francesc Pujol",
					AuthorImagePath:    "../images/cesc.jpg",
					UpdatedAt:          updatedAt,
					UpdatedBy:          sql.NullString{String: "Francesc Pujol", Valid: true},
				}

				store.EXPECT().
					UpdatePost(gomock.Any(), EqEditPostParams(arg, postId, updatedAt)).
					Times(1).
					Return(post, nil)

			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "Required Title",
			body: gin.H{
				"introduction":       "Me and you!",
				"mainImageAlt":       "Alt 1",
				"mainImagePath":      "img1",
				"thumbnailImageAlt":  "Alt 2",
				"thumbnailImagePath": "img2",
				"content":            "<p>No way</p>",
				"postCategoryID":     categoryId.String(),
				"published":          false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostParams{
					PostID:             postId,
					Title:              "Hello World",
					Slug:               "hello-world",
					Introduction:       sql.NullString{String: "Me and you!", Valid: true},
					MainImageAlt:       sql.NullString{String: "Alt 1", Valid: true},
					ThumbnailImageAlt:  sql.NullString{String: "Alt 2", Valid: true},
					Content:            sql.NullString{String: "<p>No way</p>", Valid: true},
					PostCategoryID:     categoryId,
					Author:             "Francesc Pujol",
					AuthorImagePath:    "../images/cesc.jpg",
					UpdatedAt:          updatedAt,
					UpdatedBy:          sql.NullString{String: "Francesc Pujol", Valid: true},
				}

				store.EXPECT().
					UpdatePost(gomock.Any(), EqEditPostParams(arg, postId, updatedAt)).
					Times(0)

			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Duplicate Title",
			body: gin.H{
				"postId":postId,
				"title":              "Hello World",
				"introduction":       "Me and you!",
				"mainImageAlt":       "Alt 1",
				"mainImagePath":      "img1",
				"thumbnailImageAlt":  "Alt 2",
				"thumbnailImagePath": "img2",
				"content":            "<p>No way</p>",
				"postCategoryID":     categoryId.String(),
				"published":          false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePostParams{
					PostID:             postId,
					Title:              "Hello World",
					Slug:               "hello-world",
					Introduction:       sql.NullString{String: "Me and you!", Valid: true},
					MainImageAlt:       sql.NullString{String: "Alt 1", Valid: true},
					ThumbnailImageAlt:  sql.NullString{String: "Alt 2", Valid: true},
					Content:            sql.NullString{String: "<p>No way</p>", Valid: true},
					PostCategoryID:     categoryId,
					Author:             "Francesc Pujol",
					AuthorImagePath:    "../images/cesc.jpg",
					UpdatedAt:          updatedAt,
					UpdatedBy:          sql.NullString{String: "Francesc Pujol", Valid: true},
				}

				store.EXPECT().
					UpdatePost(gomock.Any(), EqEditPostParams(arg, postId, updatedAt)).
					Times(1).
					Return(db.Post{}, &pq.Error{Code: "23505"})

			},
			want: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for i := range testTable {
		tt := testTable[i]

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			url := "/api/v1/post"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tt.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tt.want(t, recorder)
		})
	}
}