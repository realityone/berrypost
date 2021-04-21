package main

import (
	"github.com/realityone/berrypost/pkg/server"
)

func main() {
	// debug server
	server := server.New()
	server.Serve()
}
