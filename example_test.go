package socks5_test

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

func ExampleSocks5Server() {
	timeout := 60       // 60s
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
			if err := c.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
				log.Println("set local timeout:", err)
				return
			}

			s5s := &socks5.Server{
				C: c,
				SelectMethod: func(methods []byte) (method byte, got bool) {
					for _, m := range methods {
						if m == socks5.MethodNone {
							method = socks5.MethodNone
							got = true
							return
						}
					}
					return
				},
				SupportedCommands: []byte{socks5.CmdConnect},
			}
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
			if err := rc.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
				log.Println("set remote timeout:", err)
				return
			}
			go func() {
				_, _ = io.Copy(c, rc)
			}()
			_, _ = io.Copy(rc, c)

		}(c)

	}
}
