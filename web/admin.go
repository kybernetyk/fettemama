package main

import (
	"web"
	"fmt"
	"crypto/md5"
	"os"
	"time"
	"strconv"
	"mustache"
)

const (
	admin_pass = "2fe9f478faa678b1005cba27ab69c6cd"
)


var successpage = `<b>Post has been posted!</b><br><br><A href="/">Index</a>`

func checkGodLevel(ctx *web.Context) bool {
	godlevel, _ := ctx.GetSecureCookie("godlevel")
	godlevel = godHash(godlevel)
	if godlevel == admin_pass {
		return true
	}
	return false
}

func godHash(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return fmt.Sprintf("%x", hasher.Sum())
}

func createNewPost(content string) os.Error {
	post := BlogPost{
		Content:   content,
		Timestamp: time.Seconds(),
		Id:        0, //0 = create new post
	}

	_, err := Db.StorePost(&post)
	if err != nil {
		return err
	}

	return nil
}

func adminGet(ctx *web.Context) string {
	if !checkGodLevel(ctx) {
		return mustache.RenderFile("templ/admin_login.mustache")
	}

	return mustache.RenderFile("templ/admin_post.mustache")
}

func adminPost(ctx *web.Context) {
	level := ctx.Params["godlevel"]
	godlevel := godHash(level)

	if ctx.Params["what"] == "login" {
		if godlevel == admin_pass {
			ctx.SetSecureCookie("godlevel", level, 3600)
			ctx.Redirect(301, "/admin")
			return
		}
		ctx.SetSecureCookie("godlevel", "fefe", 3600)
		ctx.Redirect(301, "/")
		return
	}

	if !checkGodLevel(ctx) {
		ctx.SetSecureCookie("godlevel", "fefe", 3600)
		ctx.Redirect(301, "/")
		return
	}

	if ctx.Params["what"] == "post" {
		err := createNewPost(ctx.Params["content"])
		if err != nil {
			ctx.WriteString("couldn't post: " + err.String())
			ctx.WriteString("<br><br><A href='/'>Index</a>")
			return
		}
		ctx.WriteString(successpage)
		return
	}
}


func editGet(ctx *web.Context) string {
	if !checkGodLevel(ctx) {
		return mustache.RenderFile("templ/admin_login.mustache")
	}
	id, _ := strconv.Atoi64(ctx.Params["id"])
	post, err := Db.GetPost(id)
	if err != nil {
		return "couldn't load post with given id!"
	}
	return mustache.RenderFile("templ/admin_edit.mustache", &post)
}

func editPost(ctx *web.Context) {
	if !checkGodLevel(ctx) {
		ctx.Redirect(301, "/")
		return
	}

	id, _ := strconv.Atoi64(ctx.Params["postid"])
	post, err := Db.GetPost(id)
	if err != nil {
		ctx.WriteString("couldn't load post with given id!")
		return
	}
	post.Content = ctx.Params["content"]
	_, err = Db.StorePost(&post)
	if err != nil {
		ctx.WriteString("couldn't store post!")
		return
	}

	ctx.WriteString(successpage)
}
