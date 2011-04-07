package main
import (
	"fmt"
	//"./dumbdb"
)

type BlogPost struct {
	Content   string
	Timestamp string
	Id        int
}

type BlogRenderer interface {
	func RenderPost(post *BlogPost) string
}

func RenderPost(post_id int) string {
	post, err := g_DB.Get(post_id)
	if err != nil {
		return "error: " + err.String() + "\n"
	}
	s := fmt.Sprintf("ID: %d\nDate: %s\nContent: %s\n", post.Id, post.Timestamp, post.Content)
	return s
}
