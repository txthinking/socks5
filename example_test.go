package socks5_test

import (
	"io"
	"log"
	"net"

	"github.com/txthinking/socks5"
)

func ExampleServer() {
	socks5.Debug = true // enable socks5 debug log

	l, err := net.Listen("tcp", ":1980")
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go func(c net.Conn) {
			defer c.Close()
			s5s := socks5.NewClassicServer(c)
			if err := s5s.Negotiate(); err != nil {
				log.Println(err)
				return
			}
			r, err := s5s.GetRequest()
			if err != nil {
				log.Println(err)
				return
			}
			rc, err := r.Connect(c)
			if err != nil {
				log.Println(err)
				return
			}
			defer rc.Close()
			go func() {
				_, _ = io.Copy(c, rc)
			}()
			_, _ = io.Copy(rc, c)
		}(c)
	}
}
