package main

import (
	"flag"
	"log"

	"github.com/jmbaur/battlesnake/server"
)

func main() {
	port := flag.Int("port", 8080, "Port to run Battlesnake on")
	host := flag.String("host", "127.0.0.1", "IP or FQDN to run Battlesnake on")
	// debug := flag.Bool("debug", false, "Run the server in debug mode (verbose printing)")
	flag.Parse()

	s := server.Server{
		Host: *host,
		Port: *port,
	}
	log.Fatal(s.Run())
}
