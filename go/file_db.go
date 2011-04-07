package main

import "fmt"
import "os"
import "io/ioutil"
import "json"
import "time"

type FileDB struct {
	fetch_chan   chan fetchReq
	store_chan   chan storeReq
	update_chan  chan updateReq
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
	timestamp int64
}

type updateReq struct {
	resp_chan chan int
	err_chan  chan os.Error

	post BlogPost
}

type metaInfo struct {
	LastPostId int
	LastCommentId int
}

func getMetaInfo() metaInfo {
	fn := "posts/metainfo.json"
	contents, err := ioutil.ReadFile(fn)
	if err != nil {
		return metaInfo{0,0}
	}
	
	mi := metaInfo{}
	json.Unmarshal(contents, &mi)
	return mi
}

func saveMetaInfo(info metaInfo) {
	fn := "posts/metainfo.json"
	fmt.Println(info)
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
	db.fetch_chan = make(chan fetchReq)
	db.store_chan = make(chan storeReq)
	db.control_chan = make(chan string)
	db.update_chan = make(chan updateReq)

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

func (db *FileDB) Update(post *BlogPost) (id int, err os.Error) {
	req := updateReq{
		post: *post,
		resp_chan: make(chan int),
		err_chan: make(chan os.Error),
	}
	
	db.update_chan <- req
	
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

//post *BlogPost
func (db *FileDB) Put(content string) (id int, err os.Error) {
	req := storeReq{
		content:   content,
		timestamp:      time.Seconds(),
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
	mi := getMetaInfo()
	mi.LastPostId++;
	id = mi.LastPostId
	fn := fmt.Sprintf("posts/%d.json", id)
	post.Id = id

	bytes, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return
	}

	fmt.Println("saving post:", string(bytes))
	err = ioutil.WriteFile(fn, bytes, 0666)
	if err != nil {
		return
	}
	
	saveMetaInfo(mi)

	return
}

func updatePost(post BlogPost) (id int, err os.Error) {
	fn := fmt.Sprintf("posts/%d.json", post.Id)
	fmt.Println("updating post ...")
	bytes, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return
	}
	id = post.Id
	fmt.Println("updated post ...")
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
				Content:   req.content,
				Timestamp: req.timestamp,
			}

			id, err := savePost(post)
			if err != nil {
				req.err_chan <- err
				break
			}
			req.resp_chan <- id
		
		case req := <-db.update_chan:
			id, err := updatePost(req.post)
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
	close(db.update_chan)
	close(db.control_chan)
	
}
