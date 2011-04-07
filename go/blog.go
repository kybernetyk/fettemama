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

func RenderPost(post_id int) string {
	post, err := FetchPost(post_id)
	if err != nil {
		return "error: " + err.String() + "\n"
	}
	s := fmt.Sprintf("ID: %d\nDate: %s\nContent: %s\n", post.Id, post.Timestamp, post.Content)
	return s
}
