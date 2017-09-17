package main

import yyyp_server "github.com/shellus/yyyp/server"

func main() {
	server, err := yyyp_server.New()
	if err != nil {
		panic(err)
	}
	server.Loop()
}
