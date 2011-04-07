package main
import "os"

type BlogDB interface {
//	Put(post *BlogPost) (int, os.Error)
	Put(content string) (id int, err os.Error)
	Get(post_id int) (BlogPost, os.Error)
	
	Connect()
	Disconnect()
}

