package server_test

import (
	"testing"
	"github.com/shellus/yyyp/src/yyyp/server"
	"strconv"
)

func TestSyncNatToCloud(t *testing.T) {
	for i:= 0; i < 60; i ++{
		err := server.SyncNatToCloud("127.0.0." + strconv.Itoa(i % 10), "testing3")
		if err != nil {
			t.Error(err)
		}
	}
}