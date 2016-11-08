package socks5

import "net"

// return remote conn which u want to connect
func (r *Request) Connect(c net.Conn) (net.Conn, error) {
	rc, err := net.Dial("tcp", r.Address())
	if err != nil {
		p := NewReply(REP_HOST_UNREACHABLE, ATYP_IPV4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		if err := p.WriteTo(c); err != nil {
			return nil, err
		}
		return nil, err
	}

	a, addr, port := ParseAddress(rc.LocalAddr())
	p := NewReply(REP_SUCCESS, a, addr, port)
	if err := p.WriteTo(c); err != nil {
		return nil, err
	}

	return rc, nil
}
