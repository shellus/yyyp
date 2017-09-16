package client

import (
	"log"
	"github.com/shellus/yyyp/pack"
	"github.com/shellus/yyyp/conn"
	"errors"
	"fmt"
)

type YClient struct {
	yconn *conn.YConn
}

func (t *YClient) RequestLink(name string){
	t.yconn.WriteMessage(pack.PackLink{Name:name})
	log.Printf("send link request : [%s]", name)
}

func (t *YClient) Loop() error {

	for {
		recvPackInterface, err := t.yconn.WaitPack()
		if err != nil {
			return err
		}
		switch recvPack := recvPackInterface.(type) {

		case pack.PackPing:
			log.Printf("receive ping : [%s]", t.yconn.NetConn.RemoteAddr())

		case pack.PackPone:
			log.Printf("receive pone : [%s]", t.yconn.NetConn.RemoteAddr())

		case pack.PackConnect:
			log.Printf("receive connect command, target : [%s]", recvPack.RemoteAddr)

		case pack.PackErr:
			log.Printf("receive error : [%s]", recvPack.Message)

		}
	}
}


func New(serverAddr string)(yclient *YClient, err error){
	yconn, err := conn.Dial(serverAddr)
	if err != nil {
		err = errors.New(fmt.Sprintf("connect server err: [%s]", err))
		return
	}
	yclient = &YClient{yconn: yconn}
	return
}