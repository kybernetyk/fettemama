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


func readline(b *bufio.Reader) (p []byte, err os.Error) {
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

func clientReader(session *BlogSession) {
	var line []byte
	br := bufio.NewReader(session.conn)

	for {
		line, _ = readline(br)
		s := string(line)
		if !session.active {
			break
		}

		session.read_chan <- s
	}
}

func clientWriter(session *BlogSession) {
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

func process(session *BlogSession, user_input string) {
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

	return

	if items[0] == "read" {
		if len(items) != 2 {
			session.write_chan <- "syntax: read <post_id>\n"
			return
		}
		id, _ := strconv.Atoi(items[1])
		post, err := FetchPost(id)
		if err != nil {
			session.write_chan <- "error: " + err.String() + "\n"
			return
		}
		s := fmt.Sprintf("ID: %d\nDate: %s\nContent: %s\n", post.Id, post.Timestamp, post.Content)
		session.write_chan <- s
		return
	}

	if items[0] == "post" {
	}

}

func renderEcho(commandline string, items []string) string {
	if len(items) < 2 {
		return "syntax: echo <shit to echo>\n"
	}
	s := strings.Join(items[1:], " ")
	s += "\n"
	return s
}

func setupCMDHandlers() {
	cmd_handlers["echo"] = CommandHandler{
		handler: renderEcho,
	}

	cmd_handlers["read"] = CommandHandler{
		handler: func(commandline string, items []string) string {
			if len(items) != 2 {
				return "syntax: read <post_id>\n"
			}
			id, _ := strconv.Atoi(items[1])
			return RenderPost(id)
		},
	}

	cmd_handlers["post"] = CommandHandler{
		handler: func(commandline string, items []string) string {
			id, err := StorePost("hallo, das ist content")
			if err != nil {
				return "error: " + err.String() + "\n"
			}
			s := fmt.Sprintf("saved post with id %d\n", id)
			return s
		},
	}

}

func inputProcessor(session *BlogSession) {
	for {
		user_input := <-session.read_chan
		if !session.active {
			break
		}
		process(session, user_input)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	session := BlogSession{}

	session.write_chan = make(chan string)
	session.read_chan = make(chan string)
	session.control_chan = make(chan string)
	session.conn = conn
	session.active = true

	go clientReader(&session)
	go clientWriter(&session)
	go inputProcessor(&session)

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

func serverFunc() {
	service := fmt.Sprintf(":%d", port)
	tcpAddr, _ := net.ResolveTCPAddr(service)
	listener, _ := net.ListenTCP("tcp4", tcpAddr)

	fmt.Println("listening on: " + tcpAddr.IP.String() + service)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func RunServer() {
	db_Start()

	setupCMDHandlers()

	go serverFunc()

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
