package main

import (
	"github.com/realityone/berrypost/pkg/proxy"
	"github.com/realityone/berrypost/pkg/server"
	"github.com/realityone/berrypost/pkg/server/management"
)

func main() {
	// debug server
	components := []server.Component{}
	components = append(components, management.New(), proxy.New())

	server := server.New(server.SetComponents(components))
	server.Serve()
}
