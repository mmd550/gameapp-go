package main

import (
	httpserver "gameapp/delivery"
)

func main() {

	server := httpserver.New()

	server.Start()
}
