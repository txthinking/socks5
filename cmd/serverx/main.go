package main

import (
    "log"
    "net"
    "crypto/tls"
    //"time"
    "io"
)

func main() {
    cer, err := tls.LoadX509KeyPair(
        "/home/tx/go/src/github.com/txthinking/socks5/cmd/cert.pem",
        "/home/tx/go/src/github.com/txthinking/socks5/cmd/key.pem")
    if err != nil {
        log.Fatal(err)
        return
    }
    config := &tls.Config{
        Certificates: []tls.Certificate{cer},
        InsecureSkipVerify: true,
        MinVersion: tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        },
    }

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
        log.Println("In", c.RemoteAddr().String())
        go func (c net.Conn){
            rc, err := tls.Dial("tcp", "g.txthinking.com:20010", config)
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

        }(c)
    }
}
