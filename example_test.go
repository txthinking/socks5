package socks5_test

import "github.com/txthinking/socks5"

func ExampleServer() {
	s, err := socks5.NewClassicServer("127.0.0.1:1081", "127.0.0.1", "", "", 60, 0, 60, 60)
	if err != nil {
		panic(err)
	}
	// You can pass in custom Handler
	s.Run(nil)
}
