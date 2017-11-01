package socks5

import (
	"testing"
)

func _TestServer(t *testing.T) {
	Debug = true // enable socks5 debug log
	s, err := NewClassicServer("192.168.1.5:1081", "a", "a", 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	s.Run(nil)
}
