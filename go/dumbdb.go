package dumbdb

import "fmt"
import "./blog"

type fetchReq struct {
	resp_chan chan blog.Post
	post_id   int
}

type storeReq struct {
	resp_chan chan int

	content string
	date    string
}

var fetch_chan chan fetchReq
var store_chan chan storeReq
var control_chan chan string


func run() {

L:
	for {
		select {
		case req := <-fetch_chan:
			fmt.Println("fetch req: ", req)
			post := blog.Post{}
			post.Content = "LOLI"
			post.Timestamp = "NOW"
			req.resp_chan <- post

		case req := <-store_chan:
			fmt.Println("store req: ", req)
			req.resp_chan <- 4

		case ctl := <-control_chan:
			fmt.Println("control chan: ", ctl)

			if ctl == "stop" {
				fmt.Println("STOOOOOOOOOOOOOOOOOOOOOOOOOOOOP")
				break L
			}
		}
	}

	close(fetch_chan)
	close(store_chan)
	close(control_chan)
}

func Start() {
	fetch_chan = make(chan fetchReq)
	store_chan = make(chan storeReq)
	control_chan = make(chan string)

	go run()
}

func Stop() {
	control_chan <- "stop"
}


func FetchPost(id int) blog.Post {
	req := fetchReq{
		post_id:   id,
		resp_chan: make(chan blog.Post),
	}

	fetch_chan <- req
	r := <-req.resp_chan
	return r
}

func StorePost(content string) int {
	req := storeReq{
		content:   content,
		date:      "now",
		resp_chan: make(chan int),
	}

	store_chan <- req
	r := <-req.resp_chan
	return r
}
