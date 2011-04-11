package main

import (
	"web"
	"fmt"
	"crypto/md5"
	"os"
	"time"
	"strconv"
	"strings"
)

const (
	admin_pass = "2fe9f478faa678b1005cba27ab69c6cd"
)

var loginpage = `
<html>
<head><title>Proove your strength ...</title></head>
<body>
<form action="/admin" method="POST">

<label for="/etc">What is your godlevel?</label>
<input id="godlevel" type="text" name="godlevel"/>
<br>
<label for="shadow">Is your godlevel legit?</label>
<input id="md5" type="text" name="md5"/>
<br>
<label for="heiratswillig">There's no winter in california!'</label>
<input id="password" type="text" name="password"/>
<br>
<label for="illegal">Please write another number</label>
<input id="unused" type="text" name="unusdd"/>
<input id="what" type="hidden" value="login" name="what">
<br>
<input type="submit" name="Submit" value="Submit"/>
</form>
</body>
</html>
`

var adminpage = `
<html>
    <head><title>Project: Spanferkel</title></head>
    <body>
    <h3>GIEF POST:</h3>
    <form action="/admin" method="POST">
    <textarea rows="8" cols="80" id="content" name="content"></textarea>
    <br>
    <input id="what" type="hidden" value="post" name="what">
    <br>
    <input type="submit" name="Submit" value="Submit"/>
    </form>
    </body>
    </html>
`

var editpage = `
<html>
    <head><title>Project: Spanferkel</title></head>
    <body>
    <h3>EDIT POST:</h3>
    <form action="/admin/edit" method="POST">
    <textarea rows="8" cols="80" id="content" name="content">$postcontent$</textarea>
    <input id="postid" type="hidden" value="$postid$" name="postid">
    <br>
    <input type="submit" name="Submit" value="Submit"/>
    </form>
    </body>
    </html>
`


var successpage = `<b>Post has been posted!</b><br><br><A href="/">Index</a>`

func checkGodLevel(ctx *web.Context) bool {
	godlevel, _ := ctx.GetSecureCookie("godlevel")
	//    godlevel := godHash(level)
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
		return loginpage
	}

	return adminpage
}

func adminPost(ctx *web.Context) {
	level := ctx.Params["godlevel"]
	fmt.Println(ctx.Params["what"])
	godlevel := godHash(level)

	if ctx.Params["what"] == "login" {
		if godlevel == admin_pass {
			ctx.SetSecureCookie("godlevel", godlevel, 3600)
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
		return loginpage
	}
	id, _ := strconv.Atoi64(ctx.Params["id"])
	post, err := Db.GetPost(id)
	if err != nil {
		return "couldn't load post with given id!"
	}
	s := strings.Replace(editpage, "$postcontent$", post.Content, -1)
	s = strings.Replace(s, "$postid$", ctx.Params["id"], -1)
	return s
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
