package main

import (
	"elleFlorio/activitypub-playground/config"
	"elleFlorio/activitypub-playground/server"
)

func main() {
	config.ReadConfig()
	server.StartServer()
}
