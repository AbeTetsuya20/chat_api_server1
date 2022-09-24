package main

import (
	"diarkis-server/server"
	"fmt"
)

func main() {
	fmt.Println("Service Start!")
	server.Server()
	fmt.Println("Finish!")
}
