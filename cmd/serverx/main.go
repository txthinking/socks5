package main

import (
	"crypto/aes"
	"crypto/rand"
	"io"
	"log"
	"net"

	"github.com/txthinking/socks5"
)

func main() {
	l, err := net.Listen("tcp", ":1090")
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
		log.Println("In", c.RemoteAddr().String())
		go func(c net.Conn) {
			rc, err := net.Dial("tcp", "127.0.0.1:20010")
			if err != nil {
				log.Println(err)
				return
			}
			defer func() {
				if err := rc.Close(); err != nil {
					log.Println(err)
				}
			}()

			iv := make([]byte, aes.BlockSize)
			if _, err = io.ReadFull(rand.Reader, iv); err != nil {
				log.Println(err)
				return
			}
			if _, err := rc.Write(iv); err != nil {
				log.Println(err)
				return
			}

			crc, err := socks5.NewCipherReadWriter(rc, []byte("txthinking"), iv)
			if err != nil {
				log.Println(err)
				return
			}

			go func() {
				if _, err := io.Copy(c, crc); err != nil {
					log.Println("copy: crc->c", err)
				}
			}()
			if _, err := io.Copy(crc, c); err != nil {
				log.Println("copy: c->crc", err)
			}

		}(c)
	}
}
