package socks5

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

// ParseAddress format net.Addr to raw address
func ParseAddress(na net.Addr) (a byte, addr []byte, port []byte) {
	ss := strings.Split(na.String(), ":")
	ip := net.ParseIP(ss[0])
	ip4 := ip.To4()
	if ip4 != nil {
		a = ATYPIPv4
		addr = []byte(ip4)
	} else {
		a = ATYPIPv6
		addr = []byte(ip)
	}
	i, _ := strconv.Atoi(ss[1])
	port = make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(i))
	return
}

// ToAddress format raw address to x.x.x.x:xx
func ToAddress(a byte, addr []byte, port []byte) string {
	var h, p string
	if a == ATYPIPv4 || a == ATYPIPv6 {
		h = net.IP(addr).String()
	}
	if a == ATYPDomain {
		if len(addr) < 1 {
			return ""
		}
		if len(addr) < int(addr[0])+1 {
			return ""
		}
		h = string(addr[1:])
	}
	p = strconv.Itoa(int(binary.BigEndian.Uint16(port)))
	return net.JoinHostPort(h, p)
}
