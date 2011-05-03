package main

type BlogPost struct {
	Content   string
	Timestamp int64
	Id        int64
	Comments  []PostComment
}

func (self BlogPost) Excerpt() string {
	l := 80
	if len(self.Content) <= 80 {
		l = len(self.Content)
	}
	ex := self.Content[0:l]
	return string(ex)
}

type PostComment struct {
	Content   string
	Author    string
	Timestamp int64
	Id        int64
	PostId    int64
}


