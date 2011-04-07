package main

//import "./blog"
import "./dumbdb"
import "fmt"
import "time"

var bla = make(chan string)

func dummy() {
	time.Sleep(4000000000.0)
	bla <- "LOL"
}

func main() {
	fmt.Println("Hai!")

	dumbdb.Start()

	post := dumbdb.FetchPost(0)
	fmt.Println(post)

	id := dumbdb.StorePost("hai")
	fmt.Println(id)

	dumbdb.Stop()

	//var i int
	//	fmt.Scanf("%d", &i)

	go dummy()
L:
	for {
		select {
		case x := <-bla:
			fmt.Println("somethin:", x)
			break L
		}

	}
}
