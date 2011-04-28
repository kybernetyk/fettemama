package main

import (
	"web"
	"time"
	"strconv"
	"fmt"
	"mustache"
)

func postForId(id int64) BlogPost {
	Db := DBGet()
	defer Db.Close()

	post, _ := Db.GetPost(id)
	return post
}

func postsForDay(date *time.Time) []BlogPost {
	Db := DBGet()
	defer Db.Close()

	posts, _ := Db.GetPostsForDate(*date)
	return posts
}

func postsForMonth(date *time.Time) []BlogPost {
	Db := DBGet()
	defer Db.Close()

	posts, _ := Db.GetPostsForMonth(*date)
	return posts
}

//get css sring from cookie and embed it into a map for our mustache templates
func getCSS(ctx *web.Context) map[string]string {
	css, ok := GetCSS(ctx)
	if !ok {
		css = ""
	}
	m := map[string]string{"CSS": css}
	return m
}

// renders / 
func index(ctx *web.Context) string {
	css, ok := ctx.Params["css"]
	if ok {
		SetCSS(ctx, css)
		ctx.Redirect(302, "/")
		return "ok"
	}
	posts := postsForMonth(time.LocalTime()) //Db.GetLastNPosts(10)
	fmt.Printf("posts: %#v\n", posts)
	//embedded struct - our mustache templates need a NumOfComments field to render
	//but we don't want to put that field into the BlogPost Struct so it won't get stored
	//into the DB
	type MyPost struct {
		BlogPost
		NumOfComments int
	}

	//posts ordered by date. this is ugly. TODO: look up if mustache hase something to handle this situation
	type Date struct {
		Date  string
		Posts []MyPost
	}
	Db := DBGet()
	defer Db.Close()


	//loop through our posts and put them into the appropriate date structure
	dates := []Date{}
	var cur_date time.Time
	var date *Date
	for _, p := range posts {
		post_date := time.SecondsToLocalTime(p.Timestamp)
		if !(cur_date.Day == post_date.Day && cur_date.Month == post_date.Month && cur_date.Year == post_date.Year) {
			cur_date = *post_date
			dates = append(dates, Date{Date: cur_date.Format("Mon Jan _2 2006")})
			date = &dates[len(dates)-1]
		}
		p.Comments, _ = Db.GetComments(p.Id)
		mp := MyPost{p, len(p.Comments)}
		date.Posts = append(date.Posts, mp)
	}
	m := map[string]interface{}{
		"Dates": dates,
	}

	tmpl, _ := mustache.ParseFile("templ/index.mustache")
	s := tmpl.Render(&m, getCSS(ctx))
	return s
}

// renders /post?id=
func post(ctx *web.Context) string {
		Db := DBGet()
	defer Db.Close()

id_s := ctx.Params["id"]
	id, _ := strconv.Atoi64(id_s)
	post := postForId(id)
	post.Comments, _ = Db.GetComments(post.Id)

	type MyPost struct {
		BlogPost
		NumOfComments int
	}
	p := MyPost{post, len(post.Comments)}

	tmpl, _ := mustache.ParseFile("templ/postview.mustache")
	s := tmpl.Render(&p, getCSS(ctx))
	return s
}
