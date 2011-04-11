package main
import "os"

type BlogDB interface {
	Connect()
	Disconnect()

	StorePost(post *BlogPost) (int64, os.Error)
	GetPost(post_id int64) (BlogPost, os.Error)
	GetPostsForTimespan(start_timestamp, end_timestamp int64) (posts []BlogPost, err os.Error)
	GetLastNPosts(num_to_get int32) (posts []BlogPost, err os.Error)
	
	StoreComment(comment *PostComment) (int64, os.Error)
    GetComments(post_id int64) (comments []PostComment, err os.Error);
}

