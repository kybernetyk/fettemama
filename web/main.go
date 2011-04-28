package main

import (
	"web"
	//	"fmt"
)



func main() {
	DBConnect()
	defer DBDisconnect()

	web.Config.CookieSecret = "7C19QRmwf3mHZ9CPAaPQ0hsWeufKd"
	web.Get("/", index)
	web.Get("/post", post)

	web.Get("/rss.xml", rss)
	web.Get("/index.php/feed/", rss)
	web.Get("/index.php/feed/atom/", rss)


	web.Get("/admin/edit", editGet)
	web.Post("/admin/edit", editPost)

	web.Get("/admin", adminGet)
	web.Post("/admin", adminPost)

	web.Run("0.0.0.0:9876")

}
