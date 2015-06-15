package main

import (
	"flag"
)

func main() {

	var (
		addr string = "localhost:8081"
	)

	flag.StringVar(&addr, "addr", "localhost:8081", "")
	flag.Parse()

	server := NewServer(addr)
	StartServer(server)
}
