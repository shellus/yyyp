package server_test

import (
	"testing"
	"github.com/shellus/yyyp/src/yyyp/server"
	"strconv"
	"net"
)

func TestSyncNatToCloud(t *testing.T) {
	for i:= 0; i < 60; i ++{
		err := server.SyncNatToCloud("127.0.0." + strconv.Itoa(i % 10), "testing3")
		if err != nil {
			t.Error(err)
		}
	}
}
func TestGetCloudRecord(t *testing.T) {
	rec, err :=server.GetCloudRecord("WIN-0F1IMB3OK7F")
	if err != nil {
		return
	}
	t.Log(rec)
}
func TestUdpSend(t *testing.T) {
	sock, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return
	}
	defer sock.Close()

	rec, err :=server.GetCloudRecord("WIN-0F1IMB3OK7F")
	t.Log(rec)
	if err != nil {
		return
	}
	serverAddr, err := net.ResolveUDPAddr("udp", rec.Addr)
	if err != nil {
		return
	}
	n, err := sock.WriteTo([]byte("hello"), serverAddr)
	if err != nil {
		return
	}
	t.Log(n)

}