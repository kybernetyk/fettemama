package main

import (
	"web"
	"time"
	"strconv"
	"mustache"
	"fmt"
)

func postForId(id int64) BlogPost {
	post, _ := Db.GetPost(id)
	return post
}

func postsForDay(date *time.Time) []BlogPost {
	posts, _ := Db.GetPostsForDate(*date)
	return posts
}

func postsForMonth(date *time.Time) []BlogPost {
	posts, _ := Db.GetPostsForMonth(*date)
	return posts
}

func getCSS(ctx *web.Context) map[string]string {
	css, ok := GetCSS(ctx)
	if !ok {
		css = ""
	}
	m := map[string]string{"CSS": css}
	return m
}


func index(ctx *web.Context) string {
	css, ok := ctx.Params["css"]
	if ok {
		SetCSS(ctx, css)
		ctx.Redirect(302, "/")
		return "ok"
	}
	posts := postsForMonth(time.LocalTime()) //Db.GetLastNPosts(10)


	type MyPost struct {
		BlogPost
		NumOfComments int
	}

	type Date struct {
		Date  string
		Posts []MyPost
	}
	dates := []Date{}

	var cur_date time.Time
	date := &Date{}
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

	tmpl, _ := mustache.ParseFile("templ/all_posts.mustache")
	fmt.Printf("%#v\n", m)
	s := tmpl.Render(&m, getCSS(ctx))
	return s
}

func post(ctx *web.Context) string {
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
