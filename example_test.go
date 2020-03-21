package socks5_test

func ExampleStandardSocks5Server() {
	s, err := NewClassicServer("127.0.0.1:1080", "127.0.0.1", "", "", 60, 0, 60, 60)
	if err != nil {
		panic(err)
	}
	// You can pass in custom Handler
	s.Run(nil)
}
