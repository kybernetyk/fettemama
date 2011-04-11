package main

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


