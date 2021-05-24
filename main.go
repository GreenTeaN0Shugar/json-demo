package main

import (
	"fmt"
	"json-demo/api"
)

func main() {
	fmt.Println("IM IN MAIN")

	api.ServeHTTP()
}
