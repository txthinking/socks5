package socks5

import (
)

func (r *Request) connect(c net.Conn) error {

    rc, err := net.Dial("tcp", r.Address())
    if err != nil {
        p := NewReply(REP_HOST_UNREACHABLE, ATYP_IPV4, []byte{0,0,0,0}, []byte{0,0})
        if err := p.WriteTo(s.c); err != nil {
            log.Println(err)
        }
    }
    ss := strings.Split(rc.LocalAddr().String(), ":")

    p := NewReply(REP_SUCCESS, ATYP_IPV4, net.ParseIP(ss[0]), []byte{0,0})
    if err := p.WriteTo(s.c); err != nil {
        log.Println(err)
    }
}

