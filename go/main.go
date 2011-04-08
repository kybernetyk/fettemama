package main

func main() {
	db := NewFileDB();
	formatter := NewTelnetBlogFormatter()
	server := NewTelnetServer(db, formatter)
	
	server.Run()
}
