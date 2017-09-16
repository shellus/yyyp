package conn

import (
	"net"
	"log"
	"encoding/binary"
	"io/ioutil"
	"github.com/shellus/yyyp/pack"
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
	pingChan         chan time.Time
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
		pingChan:         make(chan time.Time),
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
		NetConn:          netConn,
		packageChan:      make(chan interface{}),
		timeoutCloseChan: make(chan bool),
		poneChan:         make(chan time.Time),
		pingChan:         make(chan time.Time),
		quit:             make(chan bool),
	}
	go yconn.waitMessage()
	go yconn.waitPing()
	return
}

func (t *YConn) Close() {
	t.quit <- true
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
func (t *YConn) waitPing() {
	for {
		select {
		case <-t.pingChan:

		case <-time.After(3 * time.Second):
			log.Printf("timeout")
			t.timeoutCloseChan <- true
		}
	}
}
func (t *YConn) loopPing() {
	for range time.Tick(time.Second) {
		t.WriteMessage(pack.PackPing{Expansion:[]byte{1}})
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


		body, err := ioutil.ReadAll(io.LimitReader(t.NetConn, int64(length)))
		if err != nil {
			log.Printf("read err: [%s]", err)
			return
		}
		log.Printf("recv body done : [%d]", len(body))
		if len(body) < 1024 {
			log.Printf("recv data [% X]", body)
		}else{
			log.Printf("recv data first 512 [% X]", body[:512])
			log.Printf("recv data last 512 [% X]", body[len(body)-512:])
		}


		recvPackInterface, err := pack.Parse(body)
		if err != nil {
			log.Printf("receive a invalid package : [%s]", err)
			return
		}
		log.Printf("recv pack [%#v]", recvPackInterface)
		// 内部吃掉pone包
		if _, ok := recvPackInterface.(*pack.PackPone); ok {
			t.poneChan <- time.Now()
			continue
		}
		// 内部回应ping包
		if _, ok := recvPackInterface.(*pack.PackPing); ok {
			t.pingChan <- time.Now()
			t.WriteMessage(pack.PackPone{})
			continue
		}


		t.packageChan <- recvPackInterface
	}
}

func (t *YConn) WriteMessage(message interface{}) (err error) {
	body, err := pack.Package(message)
	if err != nil {
		log.Printf("encode reply connect package err : [%s]", err)
		return
	}
	packageLen := uint32(len(body))
	log.Printf("header len : [%d]", packageLen)
	binary.Write(t.NetConn, binary.BigEndian, packageLen)

	n, err := t.NetConn.Write(body)

	log.Printf("send_data len : [%d]", n)

	if err != nil {
		log.Printf("send reply connect package err : [%s]", err)
		return
	}
	if n != len(body) {
		log.Printf("send reply connect package err : [%s]", "len")
		return
	}
	if len(body) < 1024 {
		log.Printf("recv data [% X]", body)
	}else{
		log.Printf("recv data first 512 [% X]", body[:512])
		log.Printf("recv data last 512 [% X]", body[len(body)-512:])
	}
	return
}
