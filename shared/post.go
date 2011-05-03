package main

type BlogPost struct {
	Content   string
	Timestamp int64
	Id        int64
	Comments  []PostComment
}

func (self BlogPost) Excerpt() string {
	t := htmlstrip(self.Content)
	l := 80
	if len(t) <= 80 {
		l = len(t)
	}
	ex := t[0:l]
	return string(ex)
}

type PostComment struct {
	Content   string
	Author    string
	Timestamp int64
	Id        int64
	PostId    int64
}


