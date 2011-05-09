package main

import (
	"web"
	"time"
	"strconv"
	//"fmt"
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

func lastPosts(num int32) []BlogPost {
	Db := DBGet()
	defer Db.Close()

	posts, _ := Db.GetLastNPosts(num)
	return posts
}

func postsForLastNDays(num int64) []BlogPost {
	Db := DBGet()
	defer Db.Close()

	posts, _ := Db.GetPostsForLastNDays(num)
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
	//posts := postsForMonth(time.LocalTime()) //Db.GetLastNPosts(10)

	//	posts := lastPosts(0xff)
	posts := postsForLastNDays(4)
	if len(posts) <= 0 {
		posts = lastPosts(23)
	}
	//fmt.Printf("posts: %#v\n", posts)
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

//renders /month?m=<>&y=
//if m and/or y are not filled the value of Today will be used
func month(ctx *web.Context) string {
	Db := DBGet()
	defer Db.Close()

	mon := ctx.Params["m"]
	yr := ctx.Params["y"]

	d := time.LocalTime()
	if len(mon) > 0 {
		d.Month, _ = strconv.Atoi(mon)
	}
	if len(yr) > 0 {
		d.Year, _ = strconv.Atoi64(yr)
	}

	posts := postsForMonth(d) //Db.GetLastNPosts(10)

	//fmt.Printf("posts: %#v\n", posts)
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
	//create PrevMonth
	pmon := *d
	pmon.Month--
	if pmon.Month <= 0 {
		pmon.Year--
		pmon.Month = 12
	}

	//fill map
	m := map[string]interface{}{
		"Dates":     dates,
		"PrevMonth": pmon,
	}

	tmpl, _ := mustache.ParseFile("templ/month.mustache")
	s := tmpl.Render(&m, getCSS(ctx))
	return s
}
