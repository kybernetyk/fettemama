package main

import (
	"fmt"
	"time"
	"strings"
)

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

	content := wordwrap(post.Content, 40)
	lines := strings.Split(content, "\n", -1)
	for _, line := range lines {
	    line = htmlstrip(line)
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
    return fmt.Sprintf("\t*[%s] %s\n", htmlstrip(comment.Author), htmlstrip(comment.Content))
}
