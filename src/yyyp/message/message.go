package message

import (
	"encoding/gob"
	"bytes"
	"github.com/shellus/pkg/logs"
)

type Action string

type Message struct {
	Action string
	data   []messageInterface
}

type messageInterface interface {
	gobEncode() []byte
}

type registerTunnelMessage struct {
	tunnelName string
}

var (
	registerTunnel Action = "register_tunnel"

)

func (s *registerTunnelMessage)gobEncode() []byte {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s)
	if err != nil {
		logs.Fatal(err)
	}
	return buf.Bytes()
}