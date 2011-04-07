package main

import (
	"fmt"
)

type BlogPost struct {
	Content   string
	Timestamp int64
	Id        int
	Comments  []PostComment
}

type PostComment struct {
	Content   string
	Author    string
	Timestamp int64
	Id        int
}

type BlogRenderer interface {
	RenderPost(post *BlogPost) string
}

type TelnetBlogRenderer struct {
	//empty for now
}

func NewTelnetBlogRenderer() *TelnetBlogRenderer {
	return &TelnetBlogRenderer{}
}

func (br *TelnetBlogRenderer) RenderPost(post *BlogPost) string {
	s := fmt.Sprintf("ID: %d\nDate: %s\nContent: %s\n", post.Id, post.Timestamp, post.Content)
	return s
}
