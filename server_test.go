package socks5

import (
	"testing"
)

func _TestServer(t *testing.T) {
	Debug = true // enable socks5 debug log
	s, err := NewClassicServer("127.0.0.1:1081", "127.0.0.1", "", "", 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	s.Run(nil)
}
