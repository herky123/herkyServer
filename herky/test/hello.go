package main

import "github.com/herky/herky/engine"

func main() {
	server := engine.NewServer("herky")
	server.Serve()
}
