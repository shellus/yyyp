package conn

import (
	"fmt"
	"os"
	"log"
)

var isDebug bool = true
var index int
var std = log.New(os.Stderr, "", log.LstdFlags)

func (t *YConn) debug(format string, v ...interface{}) {
	if isDebug {
		std.Output(2, fmt.Sprintf(fmt.Sprintf("[%d] %s", t.indexNum, format), v...))
	}
}