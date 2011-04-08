package main

import (
	"net"
	"fmt"
)


const (
	port   = 1337
)

type CommandHandler struct {
	handler        func(*BlogSession, string, []string) string
	min_perm_level int //the permission level needed to execute this command
}

var master_chan = make(chan string)
var status_chan = make(chan string)
var cmd_handlers = make(map[string]CommandHandler)

type TelnetServer struct {
	db        BlogDB
	formatter BlogFormatter
}

//// server
func NewTelnetServer(db BlogDB, formatter BlogFormatter) *TelnetServer {
	ts := TelnetServer{
		db:        db,
		formatter: formatter,
	}

	return &ts
}

func (srv *TelnetServer) setupCMDHandlers() {
    //I have no idea how to pass a reference to a function with a receiver
    //thus those ugly wrappers >.<
    //at first I had something like handler: (*BlogSession)func(string,string) string
    //in mind

	cmd_handlers["quit"] = CommandHandler{
		handler: func(session *BlogSession, commandline string, items []string) string {
			session.control_chan <- "disconnect"
			return "ok\n"
		},
		min_perm_level: 0,
	}

	cmd_handlers["die"] = CommandHandler{
		handler: func(session *BlogSession, commandline string, items []string) string {
			master_chan <- "diediedie"
			session.control_chan <- "disconnect"
			return "ok\n"
		},
		min_perm_level: 10,
	}

	cmd_handlers["auth"] = CommandHandler{
		handler: func(session *BlogSession, commandline string, items []string) string {
			return session.handleAuth(items)
		},
		min_perm_level: 0,
	}

	cmd_handlers["read"] = CommandHandler{
		handler: func(session *BlogSession, commandline string, items []string) string {
            return session.handleRead(items)
		},
		min_perm_level: 0,
	}

	cmd_handlers["post"] = CommandHandler{
		handler: func(session *BlogSession, commandline string, items []string) string {
            return session.handlePost(items)
		},
		min_perm_level: 5,
	}

	cmd_handlers["comment"] = CommandHandler{
		handler: func(session *BlogSession, commandline string, items []string) string {
            return session.handleComment(items)
		},
		min_perm_level: 0,
	}

}

//spawn new blog session for client
func (srv *TelnetServer) handleClient(conn net.Conn) {
	defer conn.Close()

	session := NewBlogSession(conn, srv.db, srv.formatter)
	go session.connReader()
	go session.connWriter()
	go session.inputProcessor()

	status_chan <- "* [" + (session.conn).RemoteAddr().String() + "] new connection"
	session.sendVersion()
	session.sendPrompt()

	for session.active {
		select {
		case status := <-session.control_chan:
			if status == "disconnect" {
				session.Disconnect()
			}
		}
	}
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
