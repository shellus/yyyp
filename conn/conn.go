package conn

import (
	"net"
	"log"
	"encoding/binary"
	"io/ioutil"
	"github.com/shellus/untitled1/pack"
	"github.com/xtaci/kcp-go"
	"fmt"
	"errors"
	"time"
	"io"
)

type YConn struct {
	NetConn          net.Conn
	packageChan      chan interface{}
	poneChan         chan time.Time
	timeoutCloseChan chan bool
	quit             chan bool
}

func Dial(serverAddr string) (yconn *YConn, err error) {
	netConn, err := kcp.Dial(serverAddr)
	if err != nil {
		return
	}
	yconn = &YConn{
		NetConn:          netConn,
		packageChan:      make(chan interface{}),
		timeoutCloseChan: make(chan bool),
		poneChan:         make(chan time.Time),
		quit:             make(chan bool),
	}
	go yconn.loopPing()
	go yconn.waitPone()
	go yconn.waitMessage()
	return
}

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
	yconn = &YConn{
		NetConn:     netConn,
		packageChan: make(chan interface{}),
	}
	go yconn.waitMessage()
	return
}

func (t *YConn) Close() {
	t.NetConn.Close()
	// todo 结束那些无限循环
}
func (t *YConn) waitPone() {
	for {
		select {
		case <-t.poneChan:

		case <-time.After(3 * time.Second):
			log.Printf("timeout")
			t.timeoutCloseChan <- true
		}
	}
}
func (t *YConn) loopPing() {
	for range time.Tick(time.Second) {
		t.WriteMessage(pack.PackPing{})
	}
}

func (t *YConn) WaitPack() (p interface{}, err error) {
	select {
	case p = <-t.packageChan:
	case <-t.timeoutCloseChan:
		err = errors.New("connect timeout")
	}
	return
}

// 接受和处理所有消息
func (t *YConn) waitMessage() {
	log.Printf("start wait message...")

	for {
		var length uint32
		err := binary.Read(t.NetConn, binary.BigEndian, &length)
		if err != nil {
			log.Printf("recv first 1 byte err : [%s]", err)
			return
		}
		log.Printf("recv header len : [%d]", length)


		data, err := ioutil.ReadAll(io.LimitReader(t.NetConn, int64(length)))
		if err != nil {
			log.Printf("read err: [%s]", err)
			return
		}
		log.Printf("recv body done : [%d]", len(data))

		recvPackInterface, err := pack.Parse(data)
		if err != nil {
			log.Printf("receive a invalid package : [%s]", err)
			return
		}

		// 内部吃掉pone包
		if _, ok := recvPackInterface.(pack.PackPone); ok {
			t.poneChan <- time.Now()
			continue
		}
		// 内部回应ping包
		if _, ok := recvPackInterface.(pack.PackPing); ok {
			t.WriteMessage(pack.PackPone{})
			continue
		}
		log.Printf("recv pack [%#v]", recvPackInterface)

		t.packageChan <- recvPackInterface
	}
}

func (t *YConn) WriteMessage(message interface{}) (err error) {
	sendData, err := pack.Package(message)
	if err != nil {
		log.Printf("encode reply connect package err : [%s]", err)
		return
	}
	packageLen := uint32(len(sendData))
	log.Printf("header len : [%d]", packageLen)
	binary.Write(t.NetConn, binary.BigEndian, packageLen)

	n, err := t.NetConn.Write(sendData)

	log.Printf("send_data len : [%d]", n)

	if err != nil {
		log.Printf("send reply connect package err : [%s]", err)
		return
	}
	if n != len(sendData) {
		log.Printf("send reply connect package err : [%s]", "len")
		return
	}
	return
}
