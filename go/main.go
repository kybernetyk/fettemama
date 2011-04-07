package main

var g_DB BlogDB

func main() {
	g_DB = NewFileDB();
	RunServer()
}
