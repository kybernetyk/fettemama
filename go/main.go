package main

func main() {
	db := NewMongoDB();
	formatter := NewTelnetBlogFormatter()
	server := NewTelnetServer(db, formatter)
	
	server.Run()
}
