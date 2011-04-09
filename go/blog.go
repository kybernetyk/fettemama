package main

import (
	"fmt"
	"time"
	"strings"
)

type BlogPost struct {
	Content   string
	Timestamp int64
	Id        int64
	Comments  []PostComment
}

type PostComment struct {
	Content   string
	Author    string
	Timestamp int64
	Id        int64
	PostId    int64
}

type BlogFormatter interface {
	FormatPost(post *BlogPost, withComments bool) string
	FormatComment(comment *PostComment) string
}

type TelnetBlogFormatter struct {
	//empty for now
}

func NewTelnetBlogFormatter() *TelnetBlogFormatter {
	return &TelnetBlogFormatter{}
}

func (bf *TelnetBlogFormatter) FormatPost(post *BlogPost, withComments bool) string {
	t := time.SecondsToLocalTime(post.Timestamp)
	s := fmt.Sprintf("Post #%d, %s\n", post.Id, t.String())
	
	lines := strings.Split(post.Content, "\n", -1)
	for _, line := range lines {
	    s += fmt.Sprintf("\t%s\n",line)
	}
	
	if !withComments{
	    return s;
	}
	    
    s += fmt.Sprintf("\nComments for post #%d:\n", post.Id)
	for _, c := range post.Comments {
        s += bf.FormatComment(&c)
	}
	return s
}

func (bf *TelnetBlogFormatter) FormatComment(comment *PostComment) string {
    return fmt.Sprintf("\t*[%s] %s\n", comment.Author, comment.Content)
}