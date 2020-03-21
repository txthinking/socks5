package socks5_test

import (
	"testing"
)

func ExampleStandardSocks5Server(t *testing.T) {
	s, err := NewClassicServer("127.0.0.1:1080", "127.0.0.1", "", "", 60, 0, 60, 60)
	if err != nil {
		panic(err)
	}
	// You can pass in custom Handler
	s.Run(nil)
}
