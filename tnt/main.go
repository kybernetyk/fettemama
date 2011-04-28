package main

func main() {
	DBConnect()
	defer DBDisconnect()

	formatter := NewTelnetBlogFormatter()
	server := NewTelnetServer(formatter)
	
	server.Run()
}
