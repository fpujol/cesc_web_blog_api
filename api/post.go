package api

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blogapi/api/request"
	"blogapi/api/response"
	db "blogapi/db/sqlc"
	"blogapi/internal/services"

	"blogapi/dtos"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (s *Server) RegisterNewPost(c *gin.Context) {
	var req request.RegisterNewPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := services.CreatePost(s.context, s.store, req)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, mapPostToResponse(result))
}

func (s *Server) EditPost(c *gin.Context) {
	var req request.EditPostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	
	result, err := services.UpdatePost(s.context, s.store, req)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, mapPostToResponse(result))
}

func (s *Server) GetPublishedPosts(c *gin.Context) {
	posts, err := s.store.ListPosts(s.context)

	if err != nil {
		fmt.Printf("error db: %v\n", err)
	}

	postsDtos := make([]dtos.IntroPost, 10)

	if posts != nil {

		for _, item := range posts {
			introPost := dtos.IntroPost{
				Title:         item.Title,
				Slug:          item.Slug,
				MainImagePath: item.MainImagePath.String,
				Introduction:  item.Introduction.String,
				AuthorName:    item.Author,
			}
			postsDtos = append(postsDtos, introPost)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": postsDtos,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello 2",
	})
}

func (s *Server) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if len(slug)==0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid slug")))
		return
	}

	result, err := s.store.GetPost(s.context, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var emptyPost db.Post
	if result == emptyPost {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, mapPostToResponse(result))
}

func (s *Server) ListPosts(c *gin.Context) {

	posts, err := s.store.ListPosts(s.context)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, mapPostsToResponse(posts))
}

func (s *Server) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file-0")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	filePath := fmt.Sprintf("http://%s/%s/%s", s.config.HTTPServerAddress, s.config.PathPostsImages, filename) //"http://localhost:5000/posts/images/" + filename

	out, err := os.Create(fmt.Sprintf("./public/%s/", s.config.PathPostsImages))

	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	res := response.UploadSunEditorResponse{Url: filePath, Name: filename}

	c.JSON(http.StatusOK, res)
}

func (s *Server) UploadMainImageFile(c *gin.Context) {	
	postId := c.Param("id")

	if len(postId)==0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid slug")))
		return
	}
	
	file, header, err := c.Request.FormFile("file-0")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	
	filePath := fmt.Sprintf("http://%s/%s/%s", s.config.HTTPServerAddress, s.config.PathPostsImages, filename) //"http://localhost:5000/posts/images/" + filename

	out, err := os.Create(fmt.Sprintf("./public/%s/%s", s.config.PathPostsImages, filename))
	
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	_, err = services.UpdateMainImagePost(s.context, s.store, postId, filePath)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := response.UploadSunEditorResponse{Url: filePath, Name: filename}

	c.JSON(http.StatusOK, res)
}

func (s *Server) UploadThumbnailImageFile(c *gin.Context) {
	postId := c.Param("id")

	if len(postId)==0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid slug")))
		return
	}

	file, header, err := c.Request.FormFile("file-0")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt

	filePath := fmt.Sprintf("http://%s/%s/%s", s.config.HTTPServerAddress, s.config.PathPostsImages, filename) //"http://localhost:5000/posts/images/" + filename
	//./public/posts/images/
	out, err := os.Create(fmt.Sprintf("./public/%s/%s", s.config.PathPostsImages, filename))
	
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	_, err = services.UpdateThubmnailImagePost(s.context, s.store, postId, filePath)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := response.UploadSunEditorResponse{Url: filePath, Name: filename}

	c.JSON(http.StatusOK, res)
}
