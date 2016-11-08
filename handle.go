package socks5

import (
	"errors"
	"net"
)

var (
	ERROR_UNSUPPORT_CMD  = errors.New("Unsupport Command")
	ERROR_USER_PASS_AUTH = errors.New("Invalid Username or Password for Auth")
)

type Server struct {
	C                 net.Conn
	CheckUserPass     func(user, pass []byte) bool
	SelectMethod      func(methods []byte) (method byte, got bool)
	SupportedCommands []byte
}

func (s *Server) Negotiate() error {
	rq, err := NewNegotiationRequestFrom(s.C)
	if err != nil {
		return err
	}
	m, got := s.SelectMethod(rq.Methods)
	if !got {
		rp := NewNegotiationReply(METHOD_UNSUPPORT_ALL)
		if err := rp.WriteTo(s.C); err != nil {
			return err
		}
	}
	rp := NewNegotiationReply(m)
	if err := rp.WriteTo(s.C); err != nil {
		return err
	}

	if m == METHOD_USERNAME_PASSWORD {
		urq, err := NewUserPassNegotiationRequestFrom(s.C)
		if err != nil {
			return err
		}
		if !s.CheckUserPass(urq.Uname, urq.Passwd) {
			urp := NewUserPassNegotiationReply(USER_PASS_STATUS_FAILURE)
			if err := urp.WriteTo(s.C); err != nil {
				return err
			}
			return ERROR_USER_PASS_AUTH
		}
		urp := NewUserPassNegotiationReply(USER_PASS_STATUS_SUCCESS)
		if err := urp.WriteTo(s.C); err != nil {
			return err
		}
	}
	return nil
}

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
		p := NewReply(REP_COMMAND_NOT_SUPPORTED, ATYP_IPV4, []byte{0, 0, 0, 0}, []byte{0, 0})
		if err := p.WriteTo(s.C); err != nil {
			return nil, err
		}
		return nil, ERROR_UNSUPPORT_CMD
	}
	return r, nil
}
