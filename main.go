package main

import (
	"flag"
)

func main() {

	var (
		addr string = ":3000"
//		addr string = ":3000"
	)

	flag.StringVar(&addr, "addr", ":3000", "")
	flag.Parse()

	server := NewServer(addr)
	StartServer(server)
}
