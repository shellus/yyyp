package conn

import (
	"net"
	"github.com/xtaci/kcp-go"
	"fmt"
	"errors"
)

type YListener struct {
	netListener net.Listener
}

func Listen(listerAddr string) (ylistener *YListener, err error) {

	listener, err := kcp.Listen(listerAddr)
	if err != nil {
		return
	}
	ylistener = &YListener{netListener: listener}
	return
}
func (t *YListener) Accept() (yconn *YConn, err error) {
	netConn, err := t.netListener.Accept()
	if err != nil {
		err = errors.New(fmt.Sprintf("server accept err : [%s]", err))
	}
	return newConn(netConn)
}