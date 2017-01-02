## socks5

[![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/socks5)](https://goreportcard.com/report/github.com/txthinking/socks5)
[![GoDoc](https://godoc.org/github.com/txthinking/socks5?status.svg)](https://godoc.org/github.com/txthinking/socks5)

SOCKS Protocol Version 5 Library

### Install
```
$ go get github.com/txthinking/socks5
```

### Example

```
func ExampleServer() {
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

```
Now you have a socks5 proxy listen on :1980
You can test with curl: `curl --socks5-hostname YOUR_SERVER_IP:1980 httpbin.org`
