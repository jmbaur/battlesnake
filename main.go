package main

import (
	"flag"
	"log"

	"github.com/jmbaur/battlesnake/server"
)

func main() {
  port := flag.Int("port", 8080, "Port to run Battlesnake on")
  host := flag.String("host", "127.0.0.1", "IP or FQDN to run Battlesnake on")
  flag.Parse()

	s := server.Server{
		Host: *host,
		Port: *port,
	}
	log.Fatal(s.Run())
}
