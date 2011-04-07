package dumbdb

import "fmt"
import "os"
import "io/ioutil"
import "json"
//import "strconv"
import "./blog"


type fetchReq struct {
	resp_chan chan blog.Post
	err_chan  chan os.Error
	post_id   int
}

type storeReq struct {
	resp_chan chan int
	err_chan  chan os.Error

	content string
	date    string
}

var fetch_chan chan fetchReq
var store_chan chan storeReq
var control_chan chan string

func openPostByID(post_id int) (post blog.Post, err os.Error) {
	fn := fmt.Sprintf("posts/%d.json", post_id)
	contents, err := ioutil.ReadFile(fn)
	if err != nil {
		//	err = os.NewError(string("post id " + strconv.Itoa(post_id) + " doesn't exist!"))
		return
	}
	//	println(string(contents))
	//io.WriteFile("filename", contents, 0x644);

	json.Unmarshal(contents, &post)

	return
}

func savePost(post blog.Post) (id int, err os.Error) {
	id = 1
	fn := fmt.Sprintf("posts/%d.json", id)
	post.Id = id
	
	bytes, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return
	}

	fmt.Println(string(bytes))
	err = ioutil.WriteFile(fn, bytes, 0666)

	return
}

func run() {
L:
	for {
		select {
		case req := <-fetch_chan:
			post, err := openPostByID(req.post_id)
			if err != nil {
				req.err_chan <- err
				break
			}
			req.resp_chan <- post

		case req := <-store_chan:
			post := blog.Post{
				Content: req.content,
				Timestamp: req.date,
			}
			
			id, err := savePost(post)
			if err != nil {
				req.err_chan <- err
				break
			}
			req.resp_chan <- id

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


func FetchPost(id int) (post blog.Post, err os.Error) {
	req := fetchReq{
		post_id:   id,
		resp_chan: make(chan blog.Post),
		err_chan:  make(chan os.Error),
	}

	fetch_chan <- req

	select {
	case p := <-req.resp_chan:
		post = p
		return
	case e := <-req.err_chan:
		err = e
		return
	}
	return
}

func StorePost(content string) (id int, err os.Error){
	req := storeReq{
		content:   content,
		date:      "now",
		resp_chan: make(chan int),
		err_chan:  make(chan os.Error),
	}

	store_chan <- req
	
	select {
	case i := <-req.resp_chan:
		id = i
		return
	case e := <-req.err_chan:
		err = e
		return
	}
	return
}
