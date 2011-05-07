package main

import (
	"log"
	"os"
)

func main() {
	DBConnect()
	defer DBDisconnect()

	f, err := os.OpenFile("../../logs/telnet.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err == nil {
		defer f.Close()
		log.SetOutput(f)
	}
	formatter := NewTelnetBlogFormatter()
	server := NewTelnetServer(formatter)

	server.Run()
}
