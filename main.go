package main

import (
	"coin-distribution/api"
	"log"
)

func main() {
	server, err := api.NewServer("127.0.0.1", 9001)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
