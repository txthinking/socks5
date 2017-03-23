package socks5

import (
	"log"
	"net"
)

// Connect remote conn which u want to connect.
// You may should write your method instead of use this method.
func (r *Request) Connect(c net.Conn) (*net.TCPConn, error) {
	if Debug {
		log.Println("Call:", r.Address())
	}
	ta, err := net.ResolveTCPAddr("tcp", r.Address())
	if err != nil {
		return nil, err
	}
	rc, err := net.DialTCP("tcp", nil, ta)
	if err != nil {
		var p *Reply
		if r.Atyp == ATYPIPv4 || r.Atyp == ATYPDomain {
			p = NewReply(RepHostUnreachable, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = NewReply(RepHostUnreachable, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(c); err != nil {
			return nil, err
		}
		return nil, err
	}

	a, addr, port := ParseAddress(rc.LocalAddr())
	p := NewReply(RepSuccess, a, addr, port)
	if err := p.WriteTo(c); err != nil {
		return nil, err
	}

	return rc, nil
}
