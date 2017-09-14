package main

import (
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"encoding/binary"
	"github.com/shellus/untitled1/pack"
	"bufio"
	"io/ioutil"
)

var nodes = make(map[string]net.Conn)


func main() {
	prot := ":8888"
	socket, err := kcp.Listen(prot)

	if err != nil {
		log.Panicf("server listen err [%s] [%s]", prot, err)
	}
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("server accept err : [%s]", err)
			continue
		}
		go handleConn(conn)

	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	var length uint32
	err := binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		log.Printf("recv first 1 byte err : [%s]", err)
		return
	}
	data, err := ioutil.ReadAll(bufio.NewReaderSize(conn, int(length)))

	recvPackInterface, err := pack.Parse(data)

	switch recvPack := recvPackInterface.(type) {

	case pack.PackReg:
		nodes[recvPack.Name] = conn
		log.Printf("client reg : [%s] [%s]", conn.RemoteAddr(), recvPack.Name)


	case pack.PackPing:
		// todo write pone pack
		log.Printf("client ping : [%s]", conn.RemoteAddr())


	case pack.PackPone:
		log.Printf("client pone : [%s]", conn.RemoteAddr())

	case pack.PackLink:
		remoteConn := nodes[recvPack.Name]

		{
			// 回复给他目标的地址
			sendData, err := pack.Package(pack.PackConnect{RemoteAddr:remoteConn.RemoteAddr().String()})
			if err != nil {
				log.Printf("encode reply connect package err : [%s]", err)
				return
			}
			n, err := conn.Write(sendData)
			if err != nil {
				log.Printf("send reply connect package err : [%s]", err)
				return
			}
			if n != len(sendData) {
				log.Printf("send reply connect package err : [%s]", "len")
				return
			}
		}

		{
			// 告诉目标
			sendData2, err := pack.Package(pack.PackConnect{RemoteAddr:conn.RemoteAddr().String()})
			if err != nil {
				log.Printf("encode reply connect package err : [%s]", err)
				return
			}
			n, err := remoteConn.Write(sendData2)
			if err != nil {
				log.Printf("send reply connect package err : [%s]", err)
				return
			}
			if n != len(sendData2) {
				log.Printf("send reply connect package err : [%s]", "len")
				return
			}
		}


	}

}
