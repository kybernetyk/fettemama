package main
import "os"

type BlogDB interface {
	Connect()
	Disconnect()

	StorePost(post *BlogPost) (int64, os.Error)
	GetPost(post_id int64) (BlogPost, os.Error)
	
	StoreComment(comment *PostComment) (int64, os.Error)
    GetComments(post_id int64) (comments []PostComment, err os.Error);
}

