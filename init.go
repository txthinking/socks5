package socks5

import (
	"log"
)

var Debug bool

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
