package main

import "fmt"
import "os"
import "io/ioutil"
import "json"

type FileDB struct {
	fetch_chan chan fetchReq
	store_chan chan storeReq
	control_chan chan string
}

type fetchReq struct {
	resp_chan chan BlogPost
	err_chan  chan os.Error
	post_id   int
}

type storeReq struct {
	resp_chan chan int
	err_chan  chan os.Error

	content string
	date    string
}

func NewFileDB() *FileDB {
	return &FileDB{}
}

func openPostByID(post_id int) (post BlogPost, err os.Error) {
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

func savePost(post BlogPost) (id int, err os.Error) {
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

func (db *FileDB) run() {
L:
	for {
		select {
		case req := <-db.fetch_chan:
			post, err := openPostByID(req.post_id)
			if err != nil {
				req.err_chan <- err
				break
			}
			req.resp_chan <- post

		case req := <-db.store_chan:
			post := BlogPost{
				Content: req.content,
				Timestamp: req.date,
			}
			
			id, err := savePost(post)
			if err != nil {
				req.err_chan <- err
				break
			}
			req.resp_chan <- id

		case ctl := <-db.control_chan:
			fmt.Println("control chan: ", ctl)

			if ctl == "stop" {
				fmt.Println("STOOOOOOOOOOOOOOOOOOOOOOOOOOOOP")
				break L
			}
		}
	}

	close(db.fetch_chan)
	close(db.store_chan)
	close(db.control_chan)
}

func (db *FileDB) Connect() {
	db.fetch_chan = make(chan fetchReq)
	db.store_chan = make(chan storeReq)
	db.control_chan = make(chan string)

	go db.run()
}

func (db *FileDB) Disconnect() {
	fmt.Println("DB IS DISCONNECTING")
	db.control_chan <- "stop"
}


func (db *FileDB) Get(id int) (post BlogPost, err os.Error) {
	req := fetchReq{
		post_id:   id,
		resp_chan: make(chan BlogPost),
		err_chan:  make(chan os.Error),
	}

	db.fetch_chan <- req

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
//post *BlogPost
func (db *FileDB) Put(content string) (id int, err os.Error){
	req := storeReq{
		content:   content,
		date:      "now",
		resp_chan: make(chan int),
		err_chan:  make(chan os.Error),
	}

	db.store_chan <- req
	
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
