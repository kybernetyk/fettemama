package main
import "os"

type BlogDB interface {
	Connect()
	Disconnect()

//	Update(post *BlogPost) (int, os.Error)
	//Put(content string) (id int, err os.Error)
	Put(post *BlogPost) (int, os.Error)
	Get(post_id int) (BlogPost, os.Error)
	
	GetMetaInfo() MetaInfo
	SaveMetaInfo(MetaInfo)
}

