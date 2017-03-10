package socks5

import (
	"errors"
	"net"
)

var (
	// ErrUnsupportCmd is the error when got unsupport command
	ErrUnsupportCmd = errors.New("Unsupport Command")
	// ErrUserPassAuth is the error when got invalid username or password
	ErrUserPassAuth = errors.New("Invalid Username or Password for Auth")
)

// Server is socks5 server wrapper
type Server struct {
	C                 net.Conn
	CheckUserPass     func(user, pass []byte) bool
	SelectMethod      func(methods []byte) (method byte, got bool)
	SupportedCommands []byte // Now only support connect command
}

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

// Negotiate handle negotiate packet.
// This method do not handle gssapi(0x01) method now.
func (s *Server) Negotiate() error {
	rq, err := NewNegotiationRequestFrom(s.C)
	if err != nil {
		return err
	}
	m, got := s.SelectMethod(rq.Methods)
	if !got {
		rp := NewNegotiationReply(MethodUnsupportAll)
		if err := rp.WriteTo(s.C); err != nil {
			return err
		}
	}
	rp := NewNegotiationReply(m)
	if err := rp.WriteTo(s.C); err != nil {
		return err
	}

	if m == MethodUsernamePassword {
		urq, err := NewUserPassNegotiationRequestFrom(s.C)
		if err != nil {
			return err
		}
		if !s.CheckUserPass(urq.Uname, urq.Passwd) {
			urp := NewUserPassNegotiationReply(UserPassStatusFailure)
			if err := urp.WriteTo(s.C); err != nil {
				return err
			}
			return ErrUserPassAuth
		}
		urp := NewUserPassNegotiationReply(UserPassStatusSuccess)
		if err := urp.WriteTo(s.C); err != nil {
			return err
		}
	}
	return nil
}

// GetRequest get request packet from client, and check command according to SupportedCommands
func (s *Server) GetRequest() (*Request, error) {
	r, err := NewRequestFrom(s.C)
	if err != nil {
		return nil, err
	}
	var supported bool
	for _, c := range s.SupportedCommands {
		if r.Cmd == c {
			supported = true
			break
		}
	}
	if !supported {
		p := NewReply(RepCommandNotSupported, ATYPIPv4, []byte{0, 0, 0, 0}, []byte{0, 0})
		if err := p.WriteTo(s.C); err != nil {
			return nil, err
		}
		return nil, ErrUnsupportCmd
	}
	return r, nil
}
