package socks5

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
)

func (r *Request) Address() string {
	var s string
	if r.Atyp == ATYP_DOMAIN {
		s = bytes.NewBuffer(r.DstAddr[1:]).String()
	} else {
		s = net.IP(r.DstAddr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(r.DstPort)))
	return net.JoinHostPort(s, p)
}
