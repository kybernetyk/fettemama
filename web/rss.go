package main

import (
	"fmt"
	"time"
	"strings"
	"web"
)

var rss_head = `
<html>
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">

<channel>
<title>fefemama</title>
<link>http://fettemama.org</link>
<description>THE BEST BLOG IN THE UNIVERSE WRITTEN IN Go :]</description>
<language>de</language>
`
var rss_item = `
<item>
<description><![CDATA[
$descriptioncontent$
]]></description>
<title><![CDATA[
$titlecontent$...]]></title>
<link>
$linkcontent$</link>
<guid>
$guidcontent$</guid>
<pubDate>
$datecontent$</pubDate>
</item>
`
var rss_footer =`
</channel>
</rss>
`

func renderRSSHeader() string {
	return rss_head
}

func renderRSSItem(post *BlogPost) string {
	s := rss_item
	s = strings.Replace(s, "$descriptioncontent$", post.Content, -1)

	s = strings.Replace(s, "$titlecontent$", post.Content[0:32], -1)

	link := fmt.Sprintf("http://fettemama.org/post?id=%d", post.Id)
	s = strings.Replace(s, "$linkcontent$", link, -1)

	guid := fmt.Sprintf("fm.post.id.%d", post.Id)
	s = strings.Replace(s, "$guidcontent$", guid, -1)
	
	post_date := time.SecondsToLocalTime(post.Timestamp)
	date := post_date.Format(time.RFC822Z)
	s = strings.Replace(s, "$datecontent$", date, -1)
	
	return s
}

func renderRSSFooter() string {
	return rss_footer
}

func RenderRSS(posts *[]BlogPost) string {
	s := ""
	s += renderRSSHeader()
	for _, p := range *posts {
		s += renderRSSItem(&p)
	}
	s += renderRSSFooter()

	return s
}

func rss(ctx *web.Context) string {
    posts,_ := Db.GetLastNPosts(20) //postsForMonth(time.LocalTime())//
    s := RenderRSS(&posts)
    s += RenderPosts(&posts)
    s += RenderFooter()
	return s
}