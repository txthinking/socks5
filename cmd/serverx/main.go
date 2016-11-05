package main

import (
    "log"
    "net"
    //"time"
    "io"
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func (){
        log.Println(http.ListenAndServe(":1094", nil))
    }()
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
        //if err := c.SetDeadline(time.Now().Add(60 * time.Second)); err != nil {
            //log.Println(err)
            //if err = c.Close(); err != nil {
                //log.Println(err)
            //}
            //continue
        //}
        go func (){
            rc, err := net.Dial("tcp", "g.txthinking.com:20010")
            if err != nil {
                log.Println(err)
                return
            }
            defer func() {
                if err := rc.Close(); err != nil {
                    log.Println(err)
                }
            }()

            go func(){
                if _, err := io.Copy(c, rc); err != nil {
                    log.Println("copy: rc->c", err)
                }
            }()
            if _, err := io.Copy(rc, c); err != nil {
                log.Println("copy: c->rc", err)
            }

        }()
    }
}
