package main

import (
	"fmt"
	"time"
	"mustache"
	"web"
)

func rss(ctx *web.Context) string {
	posts, _ := Db.GetLastNPosts(20) //postsForMonth(time.LocalTime())//
	tmpl, _ := mustache.ParseFile("templ/rss.mustache")

	type RssItem struct {
		Title       string
		Description string
		Link        string
		Guid        string
		Date        string
	}

	var items []RssItem
	for _, post := range posts {
		post_date := time.SecondsToLocalTime(post.Timestamp)
		date := post_date.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		title := htmlstrip(post.Content)
		l := len(title)
		if l > 64 {
			l = 64
		}

		item := RssItem{
			Title:       string(title[0:l]),
			Description: post.Content,
			Link:        fmt.Sprintf("http://fettemama.org/post?id=%d", post.Id),
			Guid:        fmt.Sprintf("http://fettemama.org/post?id=%d", post.Id),
			Date:        date,
		}
		items = append(items, item)

	}

	m := map[string]interface{}{"Items": items}
	return tmpl.Render(&m)
}
