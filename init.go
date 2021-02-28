package socks5

import (
	"github.com/txthinking/x"
)

// Debug enable debug log
var Debug bool
var Dial x.Dialer = x.DefaultDial
var Resolver x.Resolver = x.DefaultResolve

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
}
