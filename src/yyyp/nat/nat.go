package nat

import (
	"github.com/shellus/yyyp/src/yyyp/libs/stun"
	"net"
	"time"
	"errors"
)

var server = "stun.l.google.com:19302"

func GetNatAddr(sock *net.UDPConn)(addr string, err error) {
	serverAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return
	}
	tid := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc}
	request, err := stun.BindRequest(tid, nil, true, false)
	if err != nil {
		return
	}

	if err = sock.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return "", err
	}

	n, err := sock.WriteTo(request, serverAddr)
	if err != nil {
		return
	}
	if n < len(request) {
		err = errors.New("Short write")
		return
	}

	buf := make([]byte, 1024)
	n, _, err = sock.ReadFromUDP(buf)
	if err != nil {
		return
	}


	packet, err := stun.ParsePacket(buf[:n], nil)
	if err != nil {
		return
	}

	if packet.Error != nil {
		err = errors.New(packet.Error.Reason)
		return
	}
	if packet.Addr == nil {
		err = errors.New("STUN server didn't provide a reflexive address")
		return
	}

	addr = packet.Addr.String()
	return
}
