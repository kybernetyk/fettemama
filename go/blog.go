package blog
import (
	"fmt"
	//"./dumbdb"
)

type Post struct {
	Content   string
	Timestamp string
	Id        int
}

func RenderPost(post_id int) string {
	post, err := dumbdb.FetchPost(post_id)
	if err != nil {
		return "error: " + err.String() + "\n"
	}
	s := fmt.Sprintf("ID: %d\nDate: %s\nContent: %s\n", post.Id, post.Timestamp, post.Content)
	return s
}
