package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
)

var (
	port   = 1337
	banner = "> ISCH FICK DEINE MUDDA\n"
)

type BlogSession struct {
	write_chan   chan string
	read_chan    chan string
	control_chan chan string
	conn         net.Conn
	active       bool
}

type CommandHandler struct {
	handler func(string, []string) string
}

var master_chan = make(chan string)
var status_chan = make(chan string)
var cmd_handlers = make(map[string]CommandHandler)

type TelnetServer struct {
	db BlogDB
	renderer BlogRenderer
}

func NewTelnetServer(db BlogDB, renderer BlogRenderer) *TelnetServer {
	ts := TelnetServer {
		db: db,
		renderer: renderer,
	}
	
	return &ts
}


func (srv *TelnetServer) readline(b *bufio.Reader) (p []byte, err os.Error) {
	if p, err = b.ReadSlice('\n'); err != nil {
		return nil, err
	}
	var i int
	for i = len(p); i > 0; i-- {
		if c := p[i-1]; c != ' ' && c != '\r' && c != '\t' && c != '\n' {
			break
		}
	}
	return p[0:i], nil
}

func (srv *TelnetServer) clientReader(session *BlogSession) {
	var line []byte
	br := bufio.NewReader(session.conn)

	for {
		line, _ = srv.readline(br)
		s := string(line)
		if !session.active {
			break
		}

		session.read_chan <- s
	}
}

func (srv *TelnetServer) clientWriter(session *BlogSession) {
	var err os.Error
	for {
		b := []byte(<-session.write_chan)
		if !session.active {
			break
		}
		_, err = (session.conn).Write(b)
		if err != nil {
			session.control_chan <- "disconnect"
		}
	}
}

func (srv *TelnetServer) process(session *BlogSession, user_input string) {
	status_chan <- "* [" + (session.conn).RemoteAddr().String() + "] user input: " + user_input
	items := strings.Split(user_input, " ", -1)

	//first special commands
	if items[0] == "quit" {
		session.control_chan <- "disconnect"
		return
	}

	if items[0] == "die" {
		master_chan <- "diediedie"
		session.control_chan <- "disconnect"
		return
	}

	//now the content fetching comments
	k, ok := cmd_handlers[items[0]]
	if !ok {
		//status_chan <- "couldn't find a command for " + items[0] + "\n"
		session.write_chan <- "command not implemented\n"
		return
	}

	session.write_chan <- k.handler(user_input, items)
}


func (srv *TelnetServer) setupCMDHandlers() {
	cmd_handlers["read"] = CommandHandler{
		handler: func(commandline string, items []string) string {
			if len(items) != 2 {
				return "syntax: read <post_id>\n"
			}
			id, _ := strconv.Atoi(items[1])
			post, err := srv.db.Get(id)
			if err != nil {
				return "error: " + err.String() + "\n"
			}
			return srv.renderer.RenderPost(&post)
		},
	}

	cmd_handlers["post"] = CommandHandler{
		handler: func(commandline string, items []string) string {
			if len(items) < 2 {
				return "syntax: post <your awesome post>\n"
			}
			content := strings.Join(items[1:], " ")
			id, err := srv.db.Put(content)
			if err != nil {
				return "error: " + err.String() + "\n"
			}
			s := fmt.Sprintf("saved post with id %d\n", id)
			return s
		},
	}
	
	cmd_handlers["comment"] = CommandHandler{
		handler: func(commandline string, items []string) string {
			if len(items) < 3 {
				return "syntax: comment <post_id> <your_nick> <your many words of comment>\n"
			}
			post_id, _ := strconv.Atoi(items[1])
			post, err := srv.db.Get(post_id)
			if err != nil {
				return "error: " + err.String() + "\n"
			}
			
			mi := getMetaInfo()
			mi.LastCommentId++;
			
			nick := items[2]
			content := strings.Join(items[3:], " ")
			comment_id := mi.LastCommentId

			comment := PostComment{
				Content: content,
				Author: nick,
				Timestamp: "now",
				Id: comment_id,
			}
			post.Comments = append(post.Comments, comment)	
			fmt.Println(post.Comments)	
			i, err := srv.db.Update(&post)
			if err != nil {
				return "error: " + err.String() + "\n"
			}
			saveMetaInfo(mi);
			
			s := fmt.Sprintf("commented post with id %d\n", i)
			return s
		},
	}

}

func (srv *TelnetServer) inputProcessor(session *BlogSession) {
	for {
		user_input := <-session.read_chan
		if !session.active {
			break
		}
		srv.process(session, user_input)
	}
}

func (srv *TelnetServer) handleClient(conn net.Conn) {
	defer conn.Close()

	session := BlogSession{}

	session.write_chan = make(chan string)
	session.read_chan = make(chan string)
	session.control_chan = make(chan string)
	session.conn = conn
	session.active = true

	go srv.clientReader(&session)
	go srv.clientWriter(&session)
	go srv.inputProcessor(&session)

	status_chan <- "* [" + (session.conn).RemoteAddr().String() + "] new connection"
	session.write_chan <- banner
	b := true
	for b {
		select {
		case status := <-session.control_chan:
			//	fmt.Println("control chan: ", status)
			if status == "disconnect" {
				session.active = false
				b = false
			}
		}
	}

	close(session.read_chan)
	close(session.write_chan)
	close(session.control_chan)

	status_chan <- "* [" + (session.conn).RemoteAddr().String() + "] disconnected"
}

func (srv *TelnetServer) serverFunc() {
	service := fmt.Sprintf(":%d", port)
	tcpAddr, _ := net.ResolveTCPAddr(service)
	listener, _ := net.ListenTCP("tcp4", tcpAddr)

	fmt.Println("listening on: " + tcpAddr.IP.String() + service)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go srv.handleClient(conn)
	}
}

func (srv *TelnetServer) Run() {
	defer srv.db.Disconnect()

	srv.db.Connect()
	srv.setupCMDHandlers()
	go srv.serverFunc()

	for {
		select {
		case command := <-master_chan:
			if command == "diediedie" {
				return
			}
		case status := <-status_chan:
			fmt.Println(status)
		}
	}
}
