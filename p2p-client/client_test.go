package p2p_client_test

import (
	"testing"
	yyyp_client "github.com/shellus/yyyp/p2p-client"
)

func TestExample(t *testing.T) {

	yclient, err := yyyp_client.NewClient()
	if err != nil {
		t.Error(err)
	}
	defer yclient.Close()
	err = yclient.Run()
	if err != nil {
		t.Fatal(err)
	}
}
