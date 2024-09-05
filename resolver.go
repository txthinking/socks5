package socks5

import (
	"context"
	"net"
	"time"
)

var DNSAddrs = []string{
	"1.1.1.1:53",
	"1.0.0.1:53",
	"8.8.8.8:53",
	"8.8.4.4:53",
}

func NewResolver() (*net.Resolver, error) {
	res := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// also can use custom dialer here
			d := net.Dialer{
				Timeout: time.Second * 1,
			}
			for _, addr := range DNSAddrs {
				conn, err := d.DialContext(ctx, "udp", addr)
				if err != nil {
					continue
				}
				return conn, err
			}
			return nil, nil
		},
	}
	return res, nil
}
