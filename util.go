package socks5

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

func ParseAddress(na net.Addr) (a byte, addr []byte, port []byte) {
	ss := strings.Split(na.String(), ":")
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
	port = make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(i))
	return
}
