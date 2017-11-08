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
package main

import "github.com/txthinking/socks5"

func main() {
	socks5.Debug = true
	s, err := socks5.NewClassicServer("127.0.0.1:1080", "127.0.0.1", "", "", 0, 0, 0, 60)
	if err != nil {
		panic(err)
	}
	if err := s.Run(nil); err != nil {
		panic(err)
	}
}
```
Test with curl: `curl -x socks5://127.0.0.1:1080 http://httpbin.org/ip`

### Users: 

 * Brook [https://github.com/txthinking/brook](https://github.com/txthinking/brook)
