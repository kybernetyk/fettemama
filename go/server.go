package main

import (
	"net"
	"fmt"
	"sync"
)


const (
	port = 1337
)


type CommandHandler struct {
	handler        func(*BlogSession, []string) string
	min_perm_level int //the permission level needed to execute this command
}

var cmd_handlers = make(map[string]CommandHandler)

type TelnetServer struct {
	db        BlogDB
	formatter BlogFormatter

	control_chan chan string
	status_chan  chan string

	sessionsMutex   sync.RWMutex
	sessions    map[int]*BlogSession
	last_session_id int
}

//// server
func NewTelnetServer(db BlogDB, formatter BlogFormatter) *TelnetServer {
	ts := &TelnetServer{
		db:        db,
		formatter: formatter,
	}

	ts.control_chan = make(chan string)
	ts.status_chan = make(chan string)
	ts.sessions = make(map[int]*BlogSession)

	return ts
}

func (srv *TelnetServer) Run() {
	defer srv.db.Disconnect()

	srv.db.Connect()
	srv.setupCMDHandlers()
	go srv.serverFunc()

	for {
		select {
		case command := <-srv.control_chan:
			if command == "diediedie" {
				return
			}
		case status := <-srv.status_chan:
			fmt.Println(status)
		}
	}
}

func (srv *TelnetServer) Shutdown() {
	srv.control_chan <- "diediedie"
}

func (srv *TelnetServer) PostStatus(status string) {
    srv.status_chan <- status
}

//broadcasts message to all connected clients
func (srv *TelnetServer) Broadcast(message string) {
    for _, s := range srv.sessions {
        var ses *BlogSession = s
        ses.Send("\n*** Broadcast: " + message)
        ses.sendPrompt()
    }
}

func (srv *TelnetServer) GetUserCount() int {
	srv.sessionsMutex.RLock()
	defer srv.sessionsMutex.RUnlock()

	return len(srv.sessions)
}

func (srv *TelnetServer) registerSession(session *BlogSession) {
	srv.sessionsMutex.Lock()
	srv.last_session_id++
	session.id = srv.last_session_id
	srv.sessions[srv.last_session_id] = session
	srv.sessionsMutex.Unlock()
}

func (srv *TelnetServer) unregisterSession(session *BlogSession) {
	srv.sessionsMutex.Lock()
    srv.sessions[session.id] = nil, false
	srv.sessionsMutex.Unlock()
}


func (srv *TelnetServer) setupCMDHandlers() {
	cmd_handlers["quit"] = CommandHandler{
		handler: func(session *BlogSession, items []string) string {
			session.Disconnect()
			return "ok\n"
		},
		min_perm_level: 0,
	}

	cmd_handlers["die"] = CommandHandler{
		handler: func(session *BlogSession, items []string) string {
			session.Disconnect()
			srv.Shutdown()
			return "ok\n"
		},
		min_perm_level: 10,
	}

	cmd_handlers["auth"] = CommandHandler{
		handler:        (*BlogSession).handleAuth,
		min_perm_level: 0,
	}

	cmd_handlers["read"] = CommandHandler{
		handler:        (*BlogSession).handleRead,
		min_perm_level: 0,
	}

	cmd_handlers["post"] = CommandHandler{
		handler:        (*BlogSession).handlePost,
		min_perm_level: 5,
	}

	cmd_handlers["comment"] = CommandHandler{
		handler:        (*BlogSession).handleComment,
		min_perm_level: 0,
	}
	
	cmd_handlers["broadcast"] = CommandHandler{
	    handler: (*BlogSession).handleBroadcast,
	    min_perm_level: 0,
	}

}

//spawn new blog session for client
func (srv *TelnetServer) handleClient(conn net.Conn) {
	defer conn.Close()

	session := NewBlogSession(srv, conn)
	go session.connReader()
	go session.connWriter()
	go session.inputProcessor()

	srv.PostStatus("* [" + (session.conn).RemoteAddr().String() + "] new connection")
	srv.registerSession(session)
	
	session.sendVersion()
	s := fmt.Sprintf("There are %d users active.\n", srv.GetUserCount())
	session.Send(s)
	session.sendPrompt()

	session.Run()
	srv.PostStatus("* [" + (session.conn).RemoteAddr().String() + "] disconnected")
	srv.unregisterSession(session)
}

func (srv *TelnetServer) serverFunc() {
	service := fmt.Sprintf(":%d", port)
	tcpAddr, _ := net.ResolveTCPAddr(service)
	listener, _ := net.ListenTCP("tcp4", tcpAddr)

	srv.PostStatus ("listening on: " + tcpAddr.IP.String() + service)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go srv.handleClient(conn)
	}
}
