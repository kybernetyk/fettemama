package main

import (
	//	"net"
	//	"os"
	"fmt"
	//	"bufio"
	"strings"
	"strconv"
	"time"
	//	"crypto/md5"
)


type BlogCommand struct {
	handler        func(BlogSession, []string) string
	min_perm_level int //the permission level needed to execute this command
}

type BlogCommandHandler interface {
	AddCommand(state int, commandString string, command BlogCommand)
	HandleCommand(session BlogSession, commandline []string) string
}

type CommandMap map[string]BlogCommand

type TelnetCommandHandler struct {
	commandsByState map[int]CommandMap
}

func NewTelnetCommandHandler() *TelnetCommandHandler {
	h := &TelnetCommandHandler{}
	h.commandsByState = make(map[int]CommandMap)

	h.commandsByState[state_reading] = make(CommandMap)
	h.commandsByState[state_posting] = make(CommandMap)

	h.setupCMDHandlers()
	return h
}

func (h *TelnetCommandHandler) AddCommand(state int, commandString string, command BlogCommand) {
	cm := h.commandsByState[state]
	cm[commandString] = command
}

func (h *TelnetCommandHandler) HandleCommand(session BlogSession, commandline []string) string {
	state := session.State()
	cmdmap := h.commandsByState[state]

	//handle normal reading mode
	k, ok := cmdmap[commandline[0]]
	if !ok {
        //if user is posting we don't want to send error messages for his input
        if session.State() != state_posting {
            return "error: command not implemented\n"    
        } else {
            return ""
        }
	}
	if session.PermissionLevel() >= k.min_perm_level {
		handler := k.handler
		return handler(session, commandline)
		//		session.Send( handler(session, items) )
	} else {
		//		session.Send("error: privileges too low\n")
		return "error: privileges too low\n"
	}
	return "\n"
}

func (h *TelnetCommandHandler) setupCMDHandlers() {

	f := func(session BlogSession, items []string) string {
		session.Disconnect()
		return "ok\n"
	}
	h.AddCommand(state_reading, "quit", BlogCommand{f, 0})

	f = func(session BlogSession, items []string) string {
		session.Disconnect()
		session.Server().Shutdown()
		return "ok\n"
	}
	h.AddCommand(state_reading, "die", BlogCommand{f, 10})

	h.AddCommand(state_reading, "auth",
		BlogCommand{
			handler:        tch_handleAuth,
			min_perm_level: 0,
		})

	h.AddCommand(state_reading, "read", BlogCommand{tch_handleRead, 0})
	h.AddCommand(state_reading, "post", BlogCommand{tch_handlePost, 5})
	h.AddCommand(state_reading, "comment", BlogCommand{tch_handleComment, 0})
	h.AddCommand(state_reading, "broadcast", BlogCommand{tch_handleBroadcast, 0})
	h.AddCommand(state_reading, "help", BlogCommand{tch_handleHelp, 0})
	h.AddCommand(state_posting, "$end", BlogCommand{tch_handlePostingEnd, 0})
	
}


func tch_handleRead(session BlogSession, items []string) string {
	if len(items) != 2 {
		return "syntax: read <post_id>\n"
	}
	id, _ := strconv.Atoi64(items[1])
	post, err := session.Db().GetPost(id)
	if err != nil {
		return "error: " + err.String() + "\n"
	}

	post.Comments, _ = session.Db().GetComments(post.Id)

	return session.BlogFormatter().FormatPost(&post, true)
}

func tch_handleAuth(session BlogSession, items []string) string {
	if len(items) != 2 {
		return "syntax: auth <password>\n"
	}
	password := items[1]
	b := session.Auth(password)

	if !b {
		return "couldn't change permission level\n"
	}
	return fmt.Sprintf("permission level %d granted\n", session.PermissionLevel())
}

func tch_handlePost(session BlogSession, items []string) string {
	if len(items) != 1 {
		return "syntax: post\n"
	}
	session.ResetInputBuffer()
	session.SetState(state_posting)
	return "enter post. enter $end to end input and save post.\n01234567890123456789012345678901234567890123456789012345678901234567890123456789\n"
}

func tch_handleComment(session BlogSession, items []string) string {
	if len(items) < 3 {
		return "syntax: comment <post_id> <your_nick> <your many words of comment>\n"
	}

    post_id, _ := strconv.Atoi64(items[1])
    nick := items[2]
    content := strings.Join(items[3:], " ")
    
    comment := PostComment{
     Content:   content,
     Author:    nick,
     Timestamp: time.Seconds(),
     PostId: post_id,
    }

    i, err := session.Db().StoreComment(&comment)
    if err != nil {
     return "error: " + err.String() + "\n"
    }
    
    s := fmt.Sprintf("commented post %d. your comment's id: %d\n",post_id, i)
    return s
}

func tch_handleBroadcast(session BlogSession, items []string) string {
	if len(items) < 2 {
		return "syntax: broadcast <your broadcast>\n"
	}

	message := strings.Join(items[1:], " ")
	message += "\n"
	session.Server().Broadcast(message)

	return "Broadcast sent\n"
}


func tch_handlePostingEnd(session BlogSession, items []string) string {
	session.SetState(state_reading)

	post := BlogPost{
		Content:   strings.Trim(session.InputBuffer(), "\n\r"),
		Timestamp: time.Seconds(),
		Id:        0,   //0 = create new post
	}

	id, err := session.Db().StorePost(&post)
	if err != nil {
	 return "error: " + err.String() + "\n"
	}

	s := fmt.Sprintf("saved post with id %d\n", id)
	return s
}

func tch_handleHelp(session BlogSession, items []string) string {
	
	s := "fettemama help\n";
	s += "help\n\t* this screen\n"
	s += "comment <post_id> <your_nick> ...\n\t* add comment to a post\n"
	s += "post\n\t* create new blog post\n"
	s += "auth <password>\n\t* change user level\n"
	s += "read <post_id>\n\t* read a post\n\n"

	return s
}
