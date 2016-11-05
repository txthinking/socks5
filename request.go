package socks5

import "encoding/binary"

func (r *Request) Address() string {
	return string(r.DstAddr) + ":" + string(binary.BigEndian.Uint16(r.DstPort))
}
