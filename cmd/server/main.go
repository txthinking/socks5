package main

import (
    "log"
    "crypto/tls"
    //"time"
    "github.com/txthinking/socks5"
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
    l, err := tls.Listen("tcp", ":20010", config)

    //l, err := net.Listen("tcp", ":20010")
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
        s := socks5.NewServer(c)
        go func (){
            if err := s.Handle(); err != nil {
                log.Println(err)
            }
        }()
    }
}
