package main

import (
	"web"
//	"fmt"
	)


func main() {

web.Get("/", index)

web.Run("0.0.0.0:8080")

}
