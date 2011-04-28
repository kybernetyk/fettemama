package main

import (
	"os"
	"fmt"
	"time"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	//"strings"
	//	"html"
)

var mgoSession *mgo.Session

type MongoDB struct {
	db      *mgo.Database
	session *mgo.Session
}

func DBGet() *MongoDB {
	d := &MongoDB{}
	d.session = mgoSession.Copy()
	tmp := d.session.DB("blog")
	d.db = &tmp
	return d
}

func DBConnect() {
	var err os.Error
	mgoSession, err = mgo.Mongo("127.0.0.1")
	if err != nil {
		fmt.Println("Couldn't connect to mongo db @ localhost")
		os.Exit(-1)
		return
	}
}

func DBDisconnect() {
	mgoSession.Close()
}

func (self *MongoDB) Close() {
	self.session.Close()
}

//warning: it will marhsall the comments list - so we need to change this
//if we enable updating/editing posts
func (md *MongoDB) StorePost(post *BlogPost) (id int64, err os.Error) {
	db := md.db
	fmt.Printf("storing post: %#v\n", *post)
	//create new post
	if post.Id == 0 {
		count, _ := db.C("posts").Count()
		count++

		id = int64(count)
		post.Id = int64(count)
		err = db.C("posts").Insert(post)
		fmt.Printf("post: %#v\n", *post)
		return
	} else { //update post
		qry := bson.M{
			"id": post.Id,
		}
		err = db.C("posts").Update(qry, post)
		if err != nil {
			return
		}
	}

	return
}


func (md *MongoDB) GetPost(post_id int64) (post BlogPost, err os.Error) {
	db := md.db
	m := bson.M{"id": post_id}
	err = db.C("posts").Find(m).One(&post)
	return
}

//returns posts for a certain date
func (md *MongoDB) GetPostsForDate(date time.Time) (posts []BlogPost, err os.Error) {
	date.Hour = 0
	date.Minute = 0
	date.Second = 0

	start := date.Seconds()
	end := start + (24 * 60 * 60)

	return md.GetPostsForTimespan(start, end, -1)
}

//returns posts for a certain month
func (md *MongoDB) GetPostsForMonth(date time.Time) (posts []BlogPost, err os.Error) {
	date.Hour = 0
	date.Minute = 0
	date.Second = 0
	date.Day = 1

	next_month := date
	next_month.Month++
	if next_month.Month > 12 {
		next_month.Month = 1
		next_month.Year++
	}

	start := date.Seconds()
	end := next_month.Seconds()

	return md.GetPostsForTimespan(start, end, -1)
}


func (md *MongoDB) GetPostsForTimespan(start_timestamp, end_timestamp, order int64) (posts []BlogPost, err os.Error) {
	db := md.db

	m := bson.M{
		"$query":   bson.M{"timestamp": bson.M{"$gte": start_timestamp, "$lt": end_timestamp}},
		"$orderby": bson.M{"timestamp": order},
	}

	iter, e := db.C("posts").Find(m).Iter()
	if e != nil {
		err = e
		return
	}

	for {
		post := BlogPost{}
		e := iter.Next(&post)
		if e != nil {
			break
		}
		fmt.Printf("lol post: %#v\n", post)
		posts = append(posts, post)
	}
	return
}

func (md *MongoDB) GetLastNPosts(num_to_get int32) (posts []BlogPost, err os.Error) {
	db := md.db

	m := bson.M{
		"$query":   bson.M{},
		"$orderby": bson.M{"timestamp": -1},
	}

	iter, e := db.C("posts").Find(m).Limit(int(num_to_get)).Iter()
	if e != nil {
		err = e
		return
	}

	for {
		post := BlogPost{}
		e := iter.Next(&post)
		if e != nil {
			break
		}
		posts = append(posts, post)
	}
	return
}

func (md *MongoDB) StoreComment(comment *PostComment) (id int64, err os.Error) {
	db := md.db

	_, err = md.GetPost(comment.PostId)
	if err != nil {
		return
	}

	content := comment.Content
	author := comment.Author
	comment.Author = author   //html.EscapeString(comment.Author)
	comment.Content = content //html.EscapeString(comment.Content)

	count, _ := db.C("comments").Count()
	count++
	id = int64(count)
	comment.Id = int64(count)

	db.C("comments").Insert(comment)

	return
}

//get comments belonging to a post
func (md *MongoDB) GetComments(post_id int64) (comments []PostComment, err os.Error) {
	db := md.db

	m := bson.M{
		"$query":   bson.M{"postid": post_id},
		"$orderby": bson.M{"timestamp": 1},
	}

	iter, e := db.C("comments").Find(m).Iter()
	if e != nil {
		err = e
		return
	}

	for {
		comment := PostComment{}
		e := iter.Next(&comment)
		if e != nil {
			break
		}
		comments = append(comments, comment)
	}
	return
}
