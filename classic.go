package socks5

import "net"

// NewClassicServer return a server which allow none method and connect command
func NewClassicServer(c net.Conn) *Server {
	return &Server{
		C: c,
		SelectMethod: func(methods []byte) (method byte, got bool) {
			for _, m := range methods {
				if m == MethodNone {
					method = MethodNone
					got = true
					return
				}
			}
			return
		},
		SupportedCommands: []byte{CmdConnect},
	}
}
