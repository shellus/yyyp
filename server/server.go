package server

import (
	"log"
	"github.com/shellus/yyyp/pack"
	"github.com/shellus/yyyp/conn"
)

type YServer struct {
	yListener *conn.YListener
	connChan  chan *conn.YConn
	quitChan  chan bool
	nodes map[string]*conn.YConn
}

var listerAddr = ":8888"

func New() (yserver *YServer, err error) {

	yListener, err := conn.Listen(listerAddr)
	if err != nil {
		log.Panicf("server listen err [%s] [%s]", listerAddr, err)
		return
	}
	yserver = &YServer{yListener: yListener}
	return
}
func (t *YServer) Loop() {
	go func() {
		for {
			conn2, err := t.yListener.Accept()
			if err != nil {
				log.Printf("server accept err : [%s]", err)
				continue
			}
			t.connChan <- conn2
		}
	}()

	for {
		select {
		case conn2 := <-t.connChan:
			go t.handleConn(conn2)
		case <-t.quitChan:
			log.Printf("receive the quitChan signal")
			break
		}
	}
	log.Printf("server accept quitChan")
}

func (t *YServer) Stop() {
	t.quitChan <- true
}

func (t *YServer) handleConn(conn *conn.YConn) {

	defer conn.Close()

	for {
		recvPackInterface, err := conn.WaitPack()
		if err != nil {
			log.Printf("wait pack err : [%s] [%s]", conn.NetConn.RemoteAddr(), err)
			return
		}
		switch recvPack := recvPackInterface.(type) {

		case pack.PackReg:
			t.nodes[recvPack.Name] = conn
			log.Printf("receive reg : [%s] [%s]", conn.NetConn.RemoteAddr(), recvPack.Name)
		case pack.PackPing:
			// todo write pone pack
			log.Printf("receive ping : [%s]", conn.NetConn.RemoteAddr())

		case pack.PackPone:
			log.Printf("receive pone : [%s]", conn.NetConn.RemoteAddr())

		case pack.PackLink:
			remoteConn := t.nodes[recvPack.Name]

			{
				// 回复给他目标的地址
				err := conn.WriteMessage(pack.PackConnect{RemoteAddr: remoteConn.NetConn.RemoteAddr().String()})
				log.Printf("send message err : [%s]", err)
			}

			{
				// 告诉目标
				err := conn.WriteMessage(pack.PackConnect{RemoteAddr: conn.NetConn.RemoteAddr().String()})
				log.Printf("send message err : [%s]", err)
			}
		case pack.PackErr:
			log.Printf("receive error : [%s]", recvPack.Message)
		}
	}

}
