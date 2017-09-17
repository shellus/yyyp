package p2p_client

import (
	"fmt"
	"os"
	"log"
)

var isDebug bool = false
var std = log.New(os.Stderr, "", log.LstdFlags)

func (t *P2PClient) debug(format string, v ...interface{}) {
	if isDebug {
		std.Output(2, fmt.Sprintf(fmt.Sprintf("[%d] %s", t.indexNum, format), v...))
	}
}