package client

import "net"

type p2pListener struct {
	p2pAddr P2PAddr

	serverAddr *Addr
	udpConn *net.UDPConn

}

// 注册P2P远程端
func Listen(localId string, serverAddr *Addr) (*p2pListener, error) {

}

// 接受P2P连接
func (l *p2pListener)Accept() (P2PConn, error){

	l.udpConn.ReadFromUDP()
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l * p2pListener)Close() error{

	return
}

// Addr returns the listener's network address.
func (l * p2pListener)Addr() P2PAddr {
	return l.p2pAddr
}