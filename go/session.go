package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"strings"
//	"time"
	"crypto/md5"
)

const (
	user_pass    = "5f4dcc3b5aa765d61d8327deb882cf99" //password :]
	blogger_pass = "8dcd694437483726d7dbdbf77a862b0f"
	admin_pass   = "0b1cf6c52b2b161c34a2b163e5e6e530"
)

const (
	state_reading = 0
	state_posting = 1
)

type BlogSession interface {
    Server() *TelnetServer
    Db() BlogDB
    BlogFormatter() BlogFormatter
    Close()
    Disconnect()
    Id() int
    SetId(id int)
    Run()
    
    Send(text string)
    SendPrompt()
    SendVersion()
    
    PermissionLevel() int
    Auth(pwd string) bool
    
    State() int
    SetState(state int)
    
    ResetInputBuffer()
    InputBuffer() string
    
}

type TelnetBlogSession struct {
	conn             net.Conn
	parent_server   *TelnetServer

	write_chan       chan string
	read_chan        chan string
	control_chan     chan string

	active           bool
	permission_level int //0 - regular visitor, 5 - blogger, 10 - superuser 
	state            int
	input_buffer     string //buffer for new posts, comments, etc. which can go over multiple lines of input
	
	id              int
	
	commandHandler BlogCommandHandler
}


//var cmd_handlers = make(map[string]CommandHandler)

//////////////// session
func NewTelnetBlogSession(server *TelnetServer, conn net.Conn) *TelnetBlogSession {
	session := &TelnetBlogSession{}
	session.parent_server = server
	session.write_chan = make(chan string)
	session.read_chan = make(chan string)
	session.control_chan = make(chan string)
	session.conn = conn
	session.active = true
	session.permission_level = 0
	session.state = state_reading
	session.commandHandler = NewTelnetCommandHandler()
	return session
}

//returns the sessions parent server
func (s *TelnetBlogSession) Server() *TelnetServer {
    return s.parent_server
}

//returns the current database
func (s *TelnetBlogSession) Db() BlogDB {
    return s.Server().db
}

//returns the current formatter
func (s *TelnetBlogSession) BlogFormatter() BlogFormatter {
    return s.Server().formatter
}

//closes channels [?]
func (s *TelnetBlogSession) Close() {
    //do I have to close channels explicitely?
    
/*	close(s.read_chan)
	close(s.write_chan)
	close(s.control_chan)*/
}

//initiates disconnect
func (s *TelnetBlogSession) Disconnect() {
	s.control_chan <- "disconnect"
}

//session mainloop
func (session *TelnetBlogSession) Run() {
    for session.active {
		select {
		case status := <-session.control_chan:
			if status == "disconnect" {
				//session.Disconnect()
				session.active = false
			}
		}
	}
}

//send text
func (s *TelnetBlogSession) Send(text string) {
	s.write_chan <- text
}

func (s *TelnetBlogSession) SendPrompt() {
	if s.state == state_reading {
		s.Send("#: ")
		return
	}
	if s.state == state_posting {
		s.Send("input >\t")
		return
	}
}
func (s *TelnetBlogSession) SendVersion() {
	s.Send("fettemama.org blog system version v0.2\n\t(c) don vito 2011\n\twritten in Go\n\tuses textfiles for data storage\n\n")
}

func (s *TelnetBlogSession) Id() int {
    return s.id
}
func (s *TelnetBlogSession) SetId(id int) {
    s.id = id
}

func (s *TelnetBlogSession) PermissionLevel() int {
    return s.permission_level
}

func (s *TelnetBlogSession) Auth(pwd string) bool {
     prev_level := s.permission_level
     hasher := md5.New()
     hasher.Write([]byte(pwd))
     h_pwd := fmt.Sprintf("%x", hasher.Sum())
    
     if h_pwd == user_pass {
         s.permission_level = 0
     }
     if h_pwd == blogger_pass {
         s.permission_level = 5
     }
     if h_pwd == admin_pass {
         s.permission_level = 10
     }
    
     if prev_level == s.permission_level {
        return false
     }

    return true
}

func (s *TelnetBlogSession) State() int {
    return s.state
}
func (s *TelnetBlogSession) SetState(state int) {
    s.state = state
}

func (s *TelnetBlogSession) InputBuffer() string {
    return s.input_buffer
}

func (s *TelnetBlogSession) ResetInputBuffer() {
    s.input_buffer = ""
}

func (s *TelnetBlogSession) readline(b *bufio.Reader) (p []byte, err os.Error) {
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

func (session *TelnetBlogSession) connReader() {
	var line []byte
	br := bufio.NewReader(session.conn)

	for {
		line, _ = session.readline(br)
		s := string(line)
		if !session.active {
			break
		}

		session.read_chan <- s
	}
}

func (session *TelnetBlogSession) connWriter() {
	var err os.Error
	for {
		b := []byte(<-session.write_chan)
		if !session.active {
			break
		}
		_, err = (session.conn).Write(b)
		if err != nil {
			session.Disconnect()
		}
	}
}


func (session *TelnetBlogSession) inputProcessor() {
	for {
		user_input := <-session.read_chan
		if !session.active {
			break
		}
		session.processInput(user_input)
		session.SendPrompt()
	}
}

func (session *TelnetBlogSession) processInput(user_input string) {
	session.Server().PostStatus("* [" + (session.conn).RemoteAddr().String() + "] user input: " + user_input)
	items := strings.Split(user_input, " ", -1)
	session.input_buffer += user_input
	session.input_buffer += "\n"
    session.Send(session.commandHandler.HandleCommand(session, items))
    //handle handle command
}

