package socks5

import (
    "log"
    "net"
    "strings"
    "io"
    "strconv"
    "encoding/binary"
    //"time"
)

func (r *Request) connect(c net.Conn) error {

    log.Println("call:", r.Address())
    rc, err := net.Dial("tcp", r.Address())
    if err != nil {
        p := NewReply(REP_HOST_UNREACHABLE, ATYP_IPV4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
        if err := p.WriteTo(c); err != nil {
            return err
        }
        return err
    }
    defer func() {
        if err := rc.Close(); err != nil {
            log.Println(err)
        }
    }()
    //if err := rc.SetDeadline(time.Now().Add(60 * time.Second)); err != nil {
        //return err
    //}

    ss := strings.Split(rc.LocalAddr().String(), ":")
    var a byte
    var addr []byte
    ip := net.ParseIP(ss[0])
    ip4 := ip.To4()
    if ip4 != nil {
        a = ATYP_IPV4
        addr = []byte(ip4)
    } else {
        a = ATYP_IPV6
        addr = []byte(ip)
    }
    i, _ := strconv.Atoi(ss[1])
    port := make([]byte, 2)
    binary.BigEndian.PutUint16(port, uint16(i))
    p := NewReply(REP_SUCCESS, a, addr, port)
    if err := p.WriteTo(c); err != nil {
        return err
    }

    go func(){
        if _, err := io.Copy(c, rc); err != nil {
            log.Println("copy: rc->c", err)
        }
    }()
    if _, err := io.Copy(rc, c); err != nil {
        log.Println("copy: c->rc", err)
    }
    return nil
}
