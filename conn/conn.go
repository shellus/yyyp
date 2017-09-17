package conn

import (
	"net"
	"encoding/binary"
	"io/ioutil"
	"github.com/xtaci/kcp-go"
	"fmt"
	"errors"
	"io"
)


type YConn struct {
	NetConn     net.Conn
	packageChan chan []byte
	quitChan    chan bool
	indexNum    int
}
func newConn(netConn net.Conn) (yconn *YConn, err error) {
	index++
	yconn = &YConn{
		NetConn:     netConn,
		packageChan: make(chan []byte),
		quitChan:    make(chan bool),
		indexNum:    index,
	}
	go yconn.waitMessage()
	return
}

func (t *YConn) Close() {
	t.NetConn.Close()
	t.quitChan <- true
	close(t.packageChan)
	close(t.quitChan)
}

func Dial(serverAddr string) (yconn *YConn, err error) {
	netConn, err := kcp.Dial(serverAddr)
	if err != nil {
		return
	}
	return newConn(netConn)
}

func (t *YConn) ReadMessage() (p []byte, err error) {
	var ok bool
	select {
	case p, ok = <-t.packageChan:
		if !ok {
			panic(io.EOF)
			return
		}
		return
	case <-t.quitChan:
		err = io.EOF
		return
	}
}

// 接受和处理所有消息
func (t *YConn) waitMessage() {
	t.debug("start wait message...")

	for {
		var length uint32
		err := binary.Read(t.NetConn, binary.BigEndian, &length)
		if err != nil {
			t.debug("recv first 1 byte err : [%s]", err)
			return
		}
		t.debug("recv header len : [%d]", length)

		body, err := ioutil.ReadAll(io.LimitReader(t.NetConn, int64(length)))
		if err != nil {
			t.debug("read err: [%s]", err)
			return
		}
		t.debug("recv body done : [%d]", len(body))
		if len(body) < 1024 {
			t.debug("recv data [% X]", body)
		} else {
			t.debug("recv data first 512 [% X]", body[:512])
			t.debug("recv data last 512 [% X]", body[len(body)-512:])
		}

		t.packageChan <- body
	}
}

func (t *YConn) WriteMessage(body []byte) (err error) {
	packageLen := uint32(len(body))
	t.debug("header len : [%d]", packageLen)
	binary.Write(t.NetConn, binary.BigEndian, packageLen)

	n, err := t.NetConn.Write(body)

	t.debug("send_data len : [%d]", n)

	if err != nil {
		err = errors.New(fmt.Sprintf("send package err : [%s]", err))
		return
	}
	if n != len(body) {
		err = errors.New(fmt.Sprintf("send package len err : [%d][%d]", n, len(body)))
		return
	}
	if len(body) < 1024 {
		t.debug("recv data [% X]", body)
	} else {
		t.debug("recv data first 512 [% X]", body[:512])
		t.debug("recv data last 512 [% X]", body[len(body)-512:])
	}
	return
}
