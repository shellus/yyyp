package server

import (
	"log"
	"github.com/shellus/yyyp/pack"
	"github.com/shellus/yyyp/conn"
	"github.com/shellus/yyyp/p2p-client"
)

type YServer struct {
	yListener *conn.YListener
	connChan  chan *p2p_client.P2PClient
	quitChan  chan bool
	regNodes  map[string]*p2p_client.P2PClient
}
var listerAddr = ":8888"

func New() (yserver *YServer, err error) {

	yListener, err := conn.Listen(listerAddr)
	if err != nil {
		log.Panicf("server listen err [%s] [%s]", listerAddr, err)
		return
	}
	yserver = &YServer{
		yListener: yListener,
		connChan: make(chan *p2p_client.P2PClient),
		quitChan: make(chan bool),
	}
	return
}
func (t *YServer) Close() {
	t.quitChan <- true
}

func (t *YServer) Loop() {
	go func() {
		for {
			conn2, err := t.yListener.Accept()
			if err != nil {
				log.Printf("server accept err : [%s]", err)
				continue
			}
			yc,err := p2p_client.NewClientWithYConn(conn2)
			if err != nil {
				log.Printf("server new client err : [%s]", err)
				continue
			}
			t.connChan <- yc
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


func (t *YServer) handleConn(yclient *p2p_client.P2PClient) {

	defer yclient.Close()
	remoteAddr := yclient.YConn.NetConn.RemoteAddr()
	for {
		recvPackInterface, err := yclient.ReadPack()
		if err != nil {
			log.Printf("wait pack err : [%s] [%s]", remoteAddr, err)
			return
		}


		switch recvPack := recvPackInterface.(type) {
		case pack.PackReg:
			t.regNodes[recvPack.Name] = yclient
			log.Printf("receive reg : [%s] [%s]", remoteAddr, recvPack.Name)
		case pack.PackLink:
			remoteConn := t.regNodes[recvPack.Name]

			{
				// 回复给他目标的地址
				err := yclient.WritePack(pack.PackConnect{RemoteAddr: remoteConn.YConn.NetConn.RemoteAddr().String()})
				log.Printf("send message err : [%s]", err)
			}

			{
				// 告诉目标
				err := yclient.WritePack(pack.PackConnect{RemoteAddr: remoteAddr.String()})
				log.Printf("send message err : [%s]", err)
			}
		}
	}

}

