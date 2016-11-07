package main

import (
	"crypto/aes"
	"io"
	"log"
	"net"

	"github.com/txthinking/socks5"
)

func main() {
	l, err := net.Listen("tcp", ":20010")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Println(err)
		}
	}()
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		iv := make([]byte, aes.BlockSize)
		if _, err = io.ReadFull(c, iv); err != nil {
			log.Println(err)
			return
		}

		cc, err := socks5.NewCipherReadWriter(c, []byte("txthinking"), iv)
		if err != nil {
			log.Println(err)
			return
		}

		s := socks5.NewServer(cc)
		go func() {
			if err := s.Handle(); err != nil {
				log.Println(err)
			}
		}()
	}
}
