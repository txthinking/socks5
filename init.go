package socks5

import (
	"context"
	"net"
	"strconv"
)

var Debug bool

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func CustomResolver(r *net.Resolver, addr string) (net.Addr, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ipAddrs, err := r.LookupIPAddr(context.Background(), host)
	if err != nil {
		return nil, err
	}

	if len(ipAddrs) == 0 {
		return nil, err
	}

	// Use the first resolved IP address
	ipAddr := ipAddrs[0]
	tcpAddr := &net.TCPAddr{
		IP: ipAddr.IP,
		Port: func() int {
			intPort, err := strconv.Atoi(port)
			if err != nil {
				p, _ := net.LookupPort("tcp", port)
				return p
			}
			return intPort
		}(),
	}
	return tcpAddr, nil

}

var Resolve func(s *Server, network string, addr string) (net.Addr, error) = func(s *Server, network string, addr string) (net.Addr, error) {
	if network == "tcp" {
		if s != nil && s.Resolver != nil {
			return CustomResolver(s.Resolver, addr)
		}
		return net.ResolveTCPAddr("tcp", addr)
	}
	return net.ResolveUDPAddr("udp", addr)
}

var DialTCP func(s *Server, network string, laddr, raddr string) (net.Conn, error) = func(s *Server, network string, laddr, raddr string) (net.Conn, error) {
	var la, ra *net.TCPAddr
	if laddr != "" {
		var err error
		la, err = net.ResolveTCPAddr(network, laddr)
		if err != nil {
			return nil, err
		}
	}
	a, err := Resolve(s, network, raddr)
	if err != nil {
		return nil, err
	}
	ra = a.(*net.TCPAddr)
	return net.DialTCP(network, la, ra)
}

var DialUDP func(s *Server, network string, laddr, raddr string) (net.Conn, error) = func(s *Server, network string, laddr, raddr string) (net.Conn, error) {
	var la, ra *net.UDPAddr
	if laddr != "" {
		var err error
		la, err = net.ResolveUDPAddr(network, laddr)
		if err != nil {
			return nil, err
		}
	}
	a, err := Resolve(s, network, raddr)
	if err != nil {
		return nil, err
	}
	ra = a.(*net.UDPAddr)
	return net.DialUDP(network, la, ra)
}
