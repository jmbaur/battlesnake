package main

import (
	"github.com/jmbaur/battlesnake/server"
	"log"
)

func main() {
	s := server.Server{
		Host: "127.0.0.1",
		Port: 8080,
	}
	log.Fatal(s.Run())
}
