package client_test

import (
	"testing"
	yyyp_client "github.com/shellus/yyyp/client"
)

func TestExample(t *testing.T) {
	yclient, err := yyyp_client.New("127.0.0.1:8888")
	if err != nil {
		t.Error(err)
	}
	//yclient.RequestLink("abc")
	err = yclient.Loop()
	t.Error(err)
}
