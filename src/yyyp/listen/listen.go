package listen

import (
	"net"
	"fmt"
	"time"
	"github.com/astaxie/beego/logs"
	"github.com/shellus/yyyp/src/yyyp/nat"
	"github.com/shellus/yyyp/src/yyyp/server"
)

func Init(){
	sock, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}
	defer sock.Close()

	for {

		natAddr, err := nat.GetNatAddr(sock)
		if err != nil {
			panic(err)
		}
		err = server.SyncNatToCloud(natAddr, "macbook")
		if err != nil {
			panic(err)
		}

		if err := sock.SetDeadline(time.Now().Add(60 * time.Second)); err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		n, remote, err := sock.ReadFromUDP(buf)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				logs.Debug("timeout")
				continue
			}
			fmt.Printf("%#v\n", err)
			panic(err)
		}
		logs.Info("recv, addr %s", remote.String())

		logs.Debug("data: %s", string(buf[:n]))
	}
}

