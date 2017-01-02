package socks5

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
)

// Address return request address like x.x.x.x:xx
func (r *Request) Address() string {
	var s string
	if r.Atyp == ATYPDomain {
		s = bytes.NewBuffer(r.DstAddr[1:]).String()
	} else {
		s = net.IP(r.DstAddr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(r.DstPort)))
	return net.JoinHostPort(s, p)
}
