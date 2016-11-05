package socks5

import (
	"log"
	"net"
)

type Server struct {
	c net.Conn
}

func NewServer(c net.Conn) *Server {
	return &Server{
		c: c,
	}
}

func (s *Server) Handle() {
	defer func() {
		if err := s.c.Close(); err != nil {
			log.Println(err)
		}
	}()

	rq, err := NewNegotiationRequestFrom(s.c)
	if err != nil {
		log.Println(err)
		return
	}
	var m byte
	var got bool
	for _, m = range rq.Methods {
		if m == METHOD_NONE || m == METHOD_USERNAME_PASSWORD {
			got = true
			break
		}
	}
	if !got {
		rp := NewNegotiationReply(METHOD_UNSUPPORT_ALL)
		if err := rp.WriteTo(s.c); err != nil {
			log.Println(err)
		}
		return
	}
	rp := NewNegotiationReply(m)
	if err := rp.WriteTo(s.c); err != nil {
		log.Println(err)
		return
	}

	if m == METHOD_USERNAME_PASSWORD {
		urq, err := NewUserPassNegotiationRequestFrom(s.c)
		if err != nil {
			log.Println(err)
			return
		}
		if string(urq.Uname) != "hello" || string(urq.Passwd) != "world" {
			urp := NewUserPassNegotiationReply(USER_PASS_STATUS_FAILURE)
			if err := urp.WriteTo(s.c); err != nil {
				log.Println(err)
			}
			return
		}
	}

	r, err := NewRequestFrom(s.c)
	if err != nil {
		log.Println(err)
		return
	}
	if r.Cmd != CMD_CONNECT { // todo: need more work
		p := NewReply(REP_COMMAND_NOT_SUPPORTED, ATYP_IPV4, []byte{0, 0, 0, 0}, []byte{0, 0})
		if err := p.WriteTo(s.c); err != nil {
			log.Println(err)
		}
		return
	}

	if r.Cmd == CMD_CONNECT {
		if err := r.connect(s.c); err != nil {
			log.Println(err)
		}
		return
	}
	if r.Cmd == CMD_BIND {
		if err := r.bind(s.c); err != nil {
			log.Println(err)
		}
		return
	}
	if r.Cmd == CMD_UDP {
		if err := r.udp(s.c); err != nil {
			log.Println(err)
		}
		return
	}
}
