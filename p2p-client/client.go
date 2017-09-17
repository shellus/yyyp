package p2p_client

import (
	"log"
	"github.com/shellus/yyyp/pack"
	"github.com/shellus/yyyp/conn"
	"errors"
	"fmt"
	"time"
)

var index int

type P2PClient struct {
	YConn            *conn.YConn
	waitPackChan     chan interface{}
	poneChan         chan time.Time
	pingChan         chan time.Time
	timeoutCloseChan chan bool
	errQuitChan      chan error
	indexNum         int
}

func NewClient() (yclient *P2PClient, err error) {
	yconn, err := conn.Dial("127.0.0.1:8888")
	if err != nil {
		return
	}
	yclient, err = newClient(yconn)
	if err != nil {
		return
	}
	yclient.clientWork()
	return
}
func NewClientWithYConn(yconn *conn.YConn) (yclient *P2PClient, err error) {
	yclient, err = newClient(yconn)
	if err != nil {
		return
	}
	yclient.serverWork()
	return
}
func newClient(yconn *conn.YConn) (yclient *P2PClient, err error) {
	index++
	yclient = &P2PClient{
		YConn:            yconn,
		waitPackChan:     make(chan interface{}),
		timeoutCloseChan: make(chan bool),
		poneChan:         make(chan time.Time),
		pingChan:         make(chan time.Time),
		indexNum:         index,
	}
	go yclient.waitPackLoop()
	return
}
func (t *P2PClient) Close() {
	t.YConn.Close()
}
func (t *P2PClient) serverWork() {
	go t.waitPing()
}

func (t *P2PClient) clientWork() {
	go t.loopPing()
	go t.waitPone()
}
func (t *P2PClient) ReadPack() (ret interface{}, err error) {
	select {
	case recvPackInterface := <-t.waitPackChan:
		ret = recvPackInterface
		return
	case <-t.timeoutCloseChan:
		err = errors.New(fmt.Sprintf("connect timeout"))
		return
	case err = <-t.errQuitChan:
		return
	}
}

func (t *P2PClient) RunClient() error {

	for {
		recvPackInterface, err := t.ReadPack()
		if err != nil {
			return err
		}
		switch recvPack := recvPackInterface.(type) {

		case *pack.PackPing:
			log.Printf("receive ping : [%s]", t.YConn.NetConn.RemoteAddr())

		case *pack.PackPone:
			t.poneChan <- time.Now()
			log.Printf("receive pone : [%s]", t.YConn.NetConn.RemoteAddr())

		case *pack.PackConnect:
			log.Printf("receive connect command, target : [%s]", recvPack.RemoteAddr)

		case *pack.PackErr:
			log.Printf("receive error : [%s]", recvPack.Message)
		}
	}
}
