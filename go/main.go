package main

func main() {
	db := NewFileDB();
	renderer := NewTelnetBlogRenderer()
	server := NewTelnetServer(db, renderer)
	
	server.Run()
}
