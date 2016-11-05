package main

import (
    "net"
    "log"
    "time"

    "github.com/txthinking/socks5"
)

func main(){
    l, err := net.Listen("tcp", "1090")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := l.Close(); err != nil {
            log.Println(err)
        }
    }()
    for{
        c, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        if err := c.SetDeadline(time.Now().Add(10*time.Second)); err != nil {
            log.Println(err)
            if err = c.Close(); err != nil {
                log.Println(err)
            }
            continue
        }
        s = socks5.NewServer(c)
        go s.Handle()
    }
}

