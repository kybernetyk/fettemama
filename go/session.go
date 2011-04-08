package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"time"
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

type BlogSession struct {
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
}


//////////////// session
func NewBlogSession(server *TelnetServer, conn net.Conn) *BlogSession {
	session := &BlogSession{}
	session.parent_server = server
	session.write_chan = make(chan string)
	session.read_chan = make(chan string)
	session.control_chan = make(chan string)
	session.conn = conn
	session.active = true
	session.permission_level = 0
	session.state = state_reading
	
	return session
}

//returns the sessions parent server
func (s *BlogSession) Server() *TelnetServer {
    return s.parent_server
}

//returns the current database
func (s *BlogSession) Db() BlogDB {
    return s.Server().db
}

//returns the current formatter
func (s *BlogSession) BlogFormatter() BlogFormatter {
    return s.Server().formatter
}

//closes channels [?]
func (s *BlogSession) Close() {
    //do I have to close channels explicitely?
    
/*	close(s.read_chan)
	close(s.write_chan)
	close(s.control_chan)*/
}

//initiates disconnect
func (s *BlogSession) Disconnect() {
	s.control_chan <- "disconnect"
}

//session mainloop
func (session *BlogSession) Run() {
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
func (s *BlogSession) Send(text string) {
	s.write_chan <- text
}

func (s *BlogSession) sendPrompt() {
	if s.state == state_reading {
		s.Send("#: ")
		return
	}
	if s.state == state_posting {
		s.Send("input >\t")
		return
	}
}

func (s *BlogSession) sendVersion() {
	s.Send("fettemama.org blog system version v0.2\n\t(c) don vito 2011\n\twritten in Go\n\tuses textfiles for data storage\n\n")
}

func (s *BlogSession) readline(b *bufio.Reader) (p []byte, err os.Error) {
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

func (session *BlogSession) connReader() {
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

func (session *BlogSession) connWriter() {
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


func (session *BlogSession) inputProcessor() {
	for {
		user_input := <-session.read_chan
		if !session.active {
			break
		}
		session.processInput(user_input)
		session.sendPrompt()
	}
}

func (session *BlogSession) processInput(user_input string) {
	session.Server().PostStatus("* [" + (session.conn).RemoteAddr().String() + "] user input: " + user_input)
	items := strings.Split(user_input, " ", -1)

	//handle multiline posting mode
	if session.state == state_posting {
	    if items[0] == "$end" {
			session.state = state_reading
            
            if len(session.input_buffer) <= 0 {
                session.Send("error: post empty?\n")
                return
            }
                
			mi := session.Db().GetMetaInfo()
			mi.LastPostId++

			post := BlogPost{
				Content:   strings.Trim(session.input_buffer, "\n\r"),
				Timestamp: time.Seconds(),
				Id:        mi.LastPostId,
			}

			id, err := session.Db().Put(&post)
			if err != nil {
				session.Send("error: " + err.String() + "\n")
			}
			session.Db().SaveMetaInfo(mi)
			s := fmt.Sprintf("saved post with id %d\n", id)
			session.Send(s)
			session.input_buffer = ""
			return
		}
		session.input_buffer += user_input
		session.input_buffer += "\n"
		return;
	}

	//handle normal reading mode
	k, ok := cmd_handlers[items[0]]
	if !ok {
		session.Send("error: command not implemented\n")
		return
	}
	if session.permission_level >= k.min_perm_level {
	    handler := k.handler
		session.Send( handler(session, items) )
	} else {
		session.Send("error: privileges too low\n")
	}

}

//session handler
func (session *BlogSession) handleRead(items []string) string {
    if len(items) != 2 {
		return "syntax: read <post_id>\n"
	}
	id, _ := strconv.Atoi(items[1])
	post, err := session.Db().Get(id)
	if err != nil {
		return "error: " + err.String() + "\n"
	}
	
	return session.BlogFormatter().FormatPost(&post, true)
}

func (s *BlogSession) handleAuth(items []string) string {
    if len(items) != 2 {
    	return "syntax: auth <password>\n"
    }
    password := items[1]

	prev_level := s.permission_level

	hasher := md5.New()
	hasher.Write([]byte(password))
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
		return "couldn't change permission level\n"
	}

	return fmt.Sprintf("permission level %d granted\n", s.permission_level)
}

func (session *BlogSession) handlePost(items []string) string {
    if len(items) != 1 {
		return "syntax: post\n"
	}
	session.state = state_posting
	return "enter post. enter $end to end input and save post.\n01234567890123456789012345678901234567890123456789012345678901234567890123456789\n"
}

func (session *BlogSession) handleComment(items []string) string {
    if len(items) < 3 {
		return "syntax: comment <post_id> <your_nick> <your many words of comment>\n"
	}
	post_id, _ := strconv.Atoi(items[1])
	post, err := session.Db().Get(post_id)
	if err != nil {
		return "error: " + err.String() + "\n"
	}

	mi := session.Db().GetMetaInfo()
	mi.LastCommentId++

	nick := items[2]
	content := strings.Join(items[3:], " ")
	comment_id := mi.LastCommentId

	comment := PostComment{
		Content:   content,
		Author:    nick,
		Timestamp: time.Seconds(),
		Id:        comment_id,
	}
	post.Comments = append(post.Comments, comment)
	fmt.Println(post.Comments)
	i, err := session.Db().Put(&post)
	if err != nil {
		return "error: " + err.String() + "\n"
	}
	session.Db().SaveMetaInfo(mi)

	s := fmt.Sprintf("commented post with id %d\n", i)
	return s
}

func (session *BlogSession) handleBroadcast(items []string) string {
    if len(items) < 2 {
        return "syntax: broadcast <your broadcast>\n"
    }
    
    message := strings.Join(items[1:], " ")
    message += "\n"
    session.Server().Broadcast(message)
    
    return "Broadcast sent\n"
}