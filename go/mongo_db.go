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

	//not sure if the driver does its own locking ...
	postmu    sync.RWMutex
	commentmu sync.RWMutex
}

func NewMongoDB() *MongoDB {
	d := &MongoDB{}
	return d
}

func (md *MongoDB) Connect() {
	var err os.Error
	md.conn, err = mongo.Connect("imac.local")
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

//warning: it will marhsall the comments list - so we need to change this
//if we enable updating/editing posts
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

func (md *MongoDB) getPostsForQuery(qryobj interface{}) (posts []BlogPost, err os.Error) {
	md.postmu.Lock()
	defer md.postmu.Unlock()

	var query mongo.BSON
	query, err = mongo.Marshal(qryobj)
	if err != nil {
		return
	}

	// count, _ := md.posts.Count(query)
	// if count == 0 {
	//  err = os.NewError("COUNT 0 Post Not Found")
	//  return
	// }

	var docs *mongo.Cursor
	docs, err = md.posts.FindAll(query)
	if err != nil {
		return
	}

	var doc mongo.BSON
	for docs.HasMore() {
		doc, err = docs.GetNext()
		if err != nil {
			return
		}
		var post BlogPost
		err = mongo.Unmarshal(doc.Bytes(), &post)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	// if len(posts) == 0 {
	//     err = os.NewError("no posts found")
	// }
	return
}


func (md *MongoDB) GetPost(post_id int64) (post BlogPost, err os.Error) {
	type q map[string]interface{}

	m := q{"id": post_id}

	//find sort example
	// m := q{
	//     "$query": q{"id": q{"$gte" : post_id}},
	//     "$orderby": q{"timestamp": -1},
	// }

	var posts []BlogPost
	posts, err = md.getPostsForQuery(m)
	if err != nil || len(posts) == 0 {
		err = os.NewError("Post not found.")
		return
	}

	post = posts[0]
	return
}

func (md *MongoDB) GetPostsForTimespan(start_timestamp, end_timestamp int64) (posts []BlogPost, err os.Error) {
	type q map[string]interface{}

	//	m := q{"id": post_id}

	m := q{
		"$query":   q{"timestamp": q{"$gte": start_timestamp, "$lt": end_timestamp}},
		"$orderby": q{"timestamp": -1},
	}

	posts, err = md.getPostsForQuery(m)
	if err != nil || len(posts) == 0 {
		err = os.NewError("Posts not found.")
		return
	}

	return
}

func (md *MongoDB) GetLastNPosts(num_to_get int) (posts []BlogPost, err os.Error) {
	m := q{
	//	"$query":   q{"timestamp": q{"$gte": start_timestamp, "$lt": end_timestamp}},
		"$orderby": q{"timestamp": -1},
	}

	//var posts []BlogPost
	posts, err = md.getPostsForQuery(m)
	if err != nil || len(posts) == 0 {
		err = os.NewError("Posts not found.")
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

//get comments belonging to a post
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
