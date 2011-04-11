package main

import (
	"web"
	//	"fmt"
)


var Db *MongoDB

func main() {
	Db = NewMongoDB()
    Db.Connect()
    web.Config.CookieSecret = "7C19QRmwf3mHZ9CPAaPQ0hsWeufKd"
	web.Get("/", index)
	web.Get("/post", post)

	web.Get("/admin/edit", editGet)
	web.Post("/admin/edit", editPost)

	web.Get("/admin", adminGet)
	web.Post("/admin", adminPost)



	web.Run("0.0.0.0:8080")

}
