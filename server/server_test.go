package server_test

import (
	yyyp_server "github.com/shellus/untitled1/server"
	"testing"
)

func TestExample(t *testing.T) {
	server, err := yyyp_server.New()
	if err != nil {
		panic(err)
	}
	server.Loop()
}
