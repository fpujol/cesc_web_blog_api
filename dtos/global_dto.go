package dtos

import "time"

type EditPost struct {
	Title        string    `json:"title"`
	Introduction string    `json:"introduction"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
	AuthorName   string    `json:"authorName"`
}

type PostHeader struct {
	Title        string    `json:"title"`
	Introduction string    `json:"introduction"`
	CreatedAt    time.Time `json:"createdAt"`
	AuthorName   string    `json:"authorName"`
}

type PostDetail struct {
	DetailType string `json:"detailType"`
	Content    string `json:"content"`
	Param1     string `json:"param1"`
	Param2     string `json:"param2"`
	Param3     string `json:"param3"`
	Param4     string `json:"param4"`
}

type IntroPost struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Introduction  string `json:"introduction"`
	MainImagePath string `json:"mainImage"`
	AuthorName    string `json:"authorName"`
	CreatedAt     string `json:"createdAt"`
}

type Comment struct {
	Content string
}

type PostSlug struct {
	PostHeader   PostHeader   `json:"postHeader"`
	PostDetails  []PostDetail `json:"postDetails"`
	PostComments []Comment    `json:"postComments"`
}
