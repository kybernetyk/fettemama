package main

import "fmt"
import "os"
import "io/ioutil"
import "json"
import "sync"

type FileDB struct {
	//chans used for communication
	get_chan   chan getReq
	put_chan   chan putReq
	control_chan chan string
	
	//let's test old school mutexing
	metaInfoMutex    sync.RWMutex
}

type getReq struct {
	resp_chan chan BlogPost
	err_chan  chan os.Error
	post_id   int
}

type putReq struct {
	resp_chan chan int
	err_chan  chan os.Error

	post BlogPost
}

type MetaInfo struct {
	LastPostId int
	LastCommentId int
}

func (db *FileDB) GetMetaInfo() MetaInfo {
	db.metaInfoMutex.RLock()
	defer db.metaInfoMutex.Unlock()
	
	fn := "posts/metainfo.json"
	contents, err := ioutil.ReadFile(fn)
	if err != nil {
		return MetaInfo{0,0}
	}
	
	mi := MetaInfo{}
	json.Unmarshal(contents, &mi)
	return mi
}

func (db *FileDB) SaveMetaInfo(info MetaInfo) {
	db.metaInfoMutex.Lock()
	defer db.metaInfoMutex.Unlock()
	
	fn := "posts/metainfo.json"

	bytes, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println("error encoding metainfo:", err.String())
		return
	}
	err = ioutil.WriteFile(fn, bytes, 0666)
	if err != nil {
		fmt.Println("error saving metainfo:", err.String())
		return
	}
	
}

func NewFileDB() *FileDB {
	return &FileDB{}
}

func (db *FileDB) Connect() {
	db.get_chan = make(chan getReq)
	db.put_chan = make(chan putReq)
	db.control_chan = make(chan string)

	go db.run()
}

func (db *FileDB) Disconnect() {
	fmt.Println("DB IS DISCONNECTING")
	db.control_chan <- "stop"
}


func (db *FileDB) Get(id int) (post BlogPost, err os.Error) {
	req := getReq{
		post_id:   id,
		resp_chan: make(chan BlogPost),
		err_chan:  make(chan os.Error),
	}

	db.get_chan <- req

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

func (db *FileDB) Put(post *BlogPost) (id int, err os.Error) {
	req := putReq{
		post: *post,
		resp_chan: make(chan int),
		err_chan: make(chan os.Error),
	}
	
	db.put_chan <- req
	
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


func openPostByID(post_id int) (post BlogPost, err os.Error) {
	fn := fmt.Sprintf("posts/%d.json", post_id)
	contents, err := ioutil.ReadFile(fn)
	if err != nil {
		//	err = os.NewError(string("post id " + strconv.Itoa(post_id) + " doesn't exist!"))
		return
	}
	
	json.Unmarshal(contents, &post)

	return
}

func putPost(post BlogPost) (id int, err os.Error) {
	fn := fmt.Sprintf("posts/%d.json", post.Id)
	bytes, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return
	}
	id = post.Id
	err = ioutil.WriteFile(fn, bytes, 0666)
	return
}	

func (db *FileDB) run() {
L:
	for {
		select {
		case req := <-db.get_chan:
			post, err := openPostByID(req.post_id)
			if err != nil {
				req.err_chan <- err
				break
			}
			req.resp_chan <- post

		case req := <-db.put_chan:
			id, err := putPost(req.post)
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

	close(db.get_chan)
	close(db.put_chan)
	close(db.control_chan)
	
}
