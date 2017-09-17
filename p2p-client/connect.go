package p2p_client

import (
	"time"
	"github.com/shellus/yyyp/pack"
	"fmt"
	"errors"
)

func (t *P2PClient) waitPackLoop() {

	// 注意，这个函数只有连接中断的情况下才能return
	// 所以，如果return，一定是后续触发Close的。

	for {
		body, err := t.YConn.ReadMessage()
		if err != nil {
			err = errors.New(fmt.Sprintf("YConn ReadMessage err : [%s]", err))
			t.errQuitChan <- err
			return
		}
		recvPackInterface, err := pack.Parse(body)
		if err != nil {
			err = errors.New(fmt.Sprintf("pack Parse err : [%s]", err))
			t.errQuitChan <- err
			return
		}

		t.debug("receive pack [%#v]", recvPackInterface)

		switch recvPack := recvPackInterface.(type) {

		case *pack.PackPing:
			t.pingChan <- time.Now()
			t.WritePack(pack.PackPone{})
		case *pack.PackPone:
			t.poneChan <- time.Now()
		case *pack.PackErr:
			err = errors.New(fmt.Sprintf("remote reply a err message: [%s]", recvPack.Message))
			t.errQuitChan <- err
			return
		default:
			t.waitPackChan <- recvPack
		}
	}
}
func (t *P2PClient) waitPing() {
	for {
		select {
		case <-t.pingChan:

		case <-time.After(3 * time.Second):
			t.debug("wait ping timeout")
			t.timeoutCloseChan <- true
		}
	}
}
func (t *P2PClient) waitPone() {
	for {
		select {
		case <-t.poneChan:

		case <-time.After(3 * time.Second):
			t.debug("wait pone timeout")
			t.timeoutCloseChan <- true
		}
	}
}

func (t *P2PClient) loopPing() {
	for range time.Tick(time.Second) {
		t.WritePack(pack.PackPing{Expansion: []byte{1}})
	}
}
