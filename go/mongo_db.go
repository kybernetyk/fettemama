package main

import (
	"os"
	"fmt"
	"github.com/mikejs/gomongo/mongo"
	"sync"
)

type MongoDB struct {
	db   *mongo.Database
	conn *mongo.Connection

	posts    *mongo.Collection
	comments *mongo.Collection
	postmu   sync.RWMutex
	commentmu   sync.RWMutex
}

func NewMongoDB() *MongoDB {
	d := &MongoDB{}
	return d
}

func (md *MongoDB) Connect() {
	var err os.Error
	md.conn, err = mongo.Connect("127.0.0.1")
	if err != nil {
		fmt.Println("Couldn't connect to mongo db @ localhost")
		os.Exit(-1)
		return
	}

	md.db = md.conn.GetDB("blog")
	md.posts = md.db.GetCollection("posts")
	md.comments = md.db.GetCollection("comments")

}

func (md *MongoDB) Disconnect() {

}

func (md *MongoDB) StorePost(post *BlogPost) (id int64, err os.Error) {
	md.postmu.Lock()
	defer md.postmu.Unlock()

	qry, _ := mongo.Marshal(map[string]string{})
	count, _ := md.posts.Count(qry)
	count++

	id = count
	post.Id = count
	doc, _ := mongo.Marshal(*post)
	fmt.Println(doc)

	md.posts.Insert(doc)

	return
}

func (md *MongoDB) GetPost(post_id int64) (post BlogPost, err os.Error) {
	md.postmu.Lock()
	defer md.postmu.Unlock()

	m := map[string]int64{"id": post_id}
	var query mongo.BSON
	query, err = mongo.Marshal(m)
	if err != nil {
		return
	}

	count, _ := md.comments.Count(query)
	if count == 0 {
	    err = os.NewError("Post Not Found")
	    return
	}

	var doc mongo.BSON
	doc, err = md.posts.FindOne(query)
	if err != nil {
		err = os.NewError("Couldn't open post")
		return
	}

	err = mongo.Unmarshal(doc.Bytes(), &post)
	if err != nil {
		return
	}
	return
}

func (md *MongoDB) StoreComment(comment *PostComment) (id int64, err os.Error) {
	md.commentmu.Lock()
	defer md.commentmu.Unlock()

	//check if post with that id exists
	_, err = md.GetPost(comment.PostId)
	if err != nil {
		//err = os.NewError("Post doesn't exist :]")
		return
	}

	qry, _ := mongo.Marshal(map[string]string{})
	count, _ := md.comments.Count(qry)
	count++
	id = count
	comment.Id = count
	doc, _ := mongo.Marshal(*comment)
	fmt.Println(doc)

	md.comments.Insert(doc)

	return
}

func (md *MongoDB) GetComments(post_id int64) (comments []PostComment, err os.Error) {
    md.commentmu.Lock()
	defer md.commentmu.Unlock()
	
	m := map[string]int64{"postid": post_id}
	var query mongo.BSON
	query, err = mongo.Marshal(m)
	if err != nil {
		return
	}

	var docs *mongo.Cursor
	docs, err = md.comments.FindAll(query)
	if err != nil {
		return
	}

	var doc mongo.BSON

	for docs.HasMore() {
		doc, err = docs.GetNext()
		if err != nil {
			return
		}
		var comment PostComment
		err = mongo.Unmarshal(doc.Bytes(), &comment)
		if err != nil {
			return
		}
		comments = append(comments, comment)
	}
	return
}
