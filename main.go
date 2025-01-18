package main

import (
	"QueraMQ/server"
	"log"
)

func main() {
	server := server.NewServer("127.0.0.1:8080")
	log.Fatal(server.Run())
}
