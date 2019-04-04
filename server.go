package socks5

import (
	"errors"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	// ErrUnsupportCmd is the error when got unsupport command
	ErrUnsupportCmd = errors.New("Unsupport Command")
	// ErrUserPassAuth is the error when got invalid username or password
	ErrUserPassAuth = errors.New("Invalid Username or Password for Auth")
)

// Server is socks5 server wrapper
type Server struct {
	UserName          string
	Password          string
	Method            byte
	SupportedCommands []byte
	TCPAddr           *net.TCPAddr
	UDPAddr           *net.UDPAddr
	ServerAddr        *net.UDPAddr
	TCPListen         *net.TCPListener
	UDPConn           *net.UDPConn
	UDPExchanges      *cache.Cache
	TCPDeadline       int
	TCPTimeout        int
	UDPDeadline       int
	UDPSessionTime    int // If client does't send address, use this fixed time
	Handle            Handler
	TCPUDPAssociate   *cache.Cache
}

// UDPExchange used to store client address and remote connection
type UDPExchange struct {
	ClientAddr *net.UDPAddr
	RemoteConn *net.UDPConn
}

// NewClassicServer return a server which allow none method
func NewClassicServer(addr, ip, username, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) (*Server, error) {
	_, p, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	saddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(ip, p))
	if err != nil {
		return nil, err
	}
	m := MethodNone
	if username != "" && password != "" {
		m = MethodUsernamePassword
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs1 := cache.New(cache.NoExpiration, cache.NoExpiration)
	s := &Server{
		Method:            m,
		UserName:          username,
		Password:          password,
		SupportedCommands: []byte{CmdConnect, CmdUDP},
		TCPAddr:           taddr,
		UDPAddr:           uaddr,
		ServerAddr:        saddr,
		UDPExchanges:      cs,
		TCPTimeout:        tcpTimeout,
		TCPDeadline:       tcpDeadline,
		UDPDeadline:       udpDeadline,
		UDPSessionTime:    udpSessionTime,
		TCPUDPAssociate:   cs1,
	}
	return s, nil
}

// Negotiate handle negotiate packet.
// This method do not handle gssapi(0x01) method now.
// Error or OK both replied.
func (s *Server) Negotiate(c *net.TCPConn) error {
	rq, err := NewNegotiationRequestFrom(c)
	if err != nil {
		return err
	}
	var got bool
	var m byte
	for _, m = range rq.Methods {
		if m == s.Method {
			got = true
		}
	}
	if !got {
		rp := NewNegotiationReply(MethodUnsupportAll)
		if err := rp.WriteTo(c); err != nil {
			return err
		}
	}
	rp := NewNegotiationReply(s.Method)
	if err := rp.WriteTo(c); err != nil {
		return err
	}

	if s.Method == MethodUsernamePassword {
		urq, err := NewUserPassNegotiationRequestFrom(c)
		if err != nil {
			return err
		}
		if string(urq.Uname) != s.UserName || string(urq.Passwd) != s.Password {
			urp := NewUserPassNegotiationReply(UserPassStatusFailure)
			if err := urp.WriteTo(c); err != nil {
				return err
			}
			return ErrUserPassAuth
		}
		urp := NewUserPassNegotiationReply(UserPassStatusSuccess)
		if err := urp.WriteTo(c); err != nil {
			return err
		}
	}
	return nil
}

// GetRequest get request packet from client, and check command according to SupportedCommands
// Error replied.
func (s *Server) GetRequest(c *net.TCPConn) (*Request, error) {
	r, err := NewRequestFrom(c)
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
		var p *Reply
		if r.Atyp == ATYPIPv4 || r.Atyp == ATYPDomain {
			p = NewReply(RepCommandNotSupported, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = NewReply(RepCommandNotSupported, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(c); err != nil {
			return nil, err
		}
		return nil, ErrUnsupportCmd
	}
	return r, nil
}

// Run server
func (s *Server) Run(h Handler) error {
	if h == nil {
		s.Handle = &DefaultHandle{}
	} else {
		s.Handle = h
	}
	errch := make(chan error)
	go func() {
		errch <- s.RunTCPServer()
	}()
	go func() {
		errch <- s.RunUDPServer()
	}()
	return <-errch
}

// RunTCPServer starts tcp server
func (s *Server) RunTCPServer() error {
	var err error
	s.TCPListen, err = net.ListenTCP("tcp", s.TCPAddr)
	if err != nil {
		return err
	}
	defer s.TCPListen.Close()
	for {
		c, err := s.TCPListen.AcceptTCP()
		if err != nil {
			return err
		}
		go func(c *net.TCPConn) {
			defer c.Close()
			if s.TCPTimeout != 0 {
				if err := c.SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
					log.Println(err)
					return
				}
			}
			if s.TCPDeadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
					log.Println(err)
					return
				}
			}
			if err := s.Negotiate(c); err != nil {
				log.Println(err)
				return
			}
			r, err := s.GetRequest(c)
			if err != nil {
				log.Println(err)
				return
			}
			if err := s.Handle.TCPHandle(s, c, r); err != nil {
				log.Println(err)
			}
		}(c)
	}
	return nil
}

// RunUDPServer starts udp server
func (s *Server) RunUDPServer() error {
	var err error
	s.UDPConn, err = net.ListenUDP("udp", s.UDPAddr)
	if err != nil {
		return err
	}
	defer s.UDPConn.Close()
	for {
		b := make([]byte, 65536)
		n, addr, err := s.UDPConn.ReadFromUDP(b)
		if err != nil {
			return err
		}
		go func(addr *net.UDPAddr, b []byte) {
			d, err := NewDatagramFromBytes(b)
			if err != nil {
				log.Println(err)
				return
			}
			if d.Frag != 0x00 {
				log.Println("Ignore frag", d.Frag)
				return
			}
			if err := s.Handle.UDPHandle(s, addr, d); err != nil {
				log.Println(err)
				return
			}
		}(addr, b[0:n])
	}
	return nil
}

// Stop server
func (s *Server) Stop() error {
	var err, err1 error
	if s.TCPListen != nil {
		err = s.TCPListen.Close()
	}
	if s.UDPConn != nil {
		err1 = s.UDPConn.Close()
	}
	if err != nil {
		return err
	}
	return err1
}

// Handler handle tcp, udp request
type Handler interface {
	// Request has not been replied yet
	TCPHandle(*Server, *net.TCPConn, *Request) error
	UDPHandle(*Server, *net.UDPAddr, *Datagram) error
}

// DefaultHandle implements Handler interface
type DefaultHandle struct {
}

// TCPHandle auto handle request. You may prefer to do yourself.
func (h *DefaultHandle) TCPHandle(s *Server, c *net.TCPConn, r *Request) error {
	if r.Cmd == CmdConnect {
		rc, err := r.Connect(c)
		if err != nil {
			return err
		}
		defer rc.Close()
		go func() {
			var bf [1024 * 2]byte
			for {
				if s.TCPDeadline != 0 {
					if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
						return
					}
				}
				i, err := rc.Read(bf[:])
				if err != nil {
					return
				}
				if _, err := c.Write(bf[0:i]); err != nil {
					return
				}
			}
		}()
		var bf [1024 * 2]byte
		for {
			if s.TCPDeadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
					return nil
				}
			}
			i, err := c.Read(bf[:])
			if err != nil {
				return nil
			}
			if _, err := rc.Write(bf[0:i]); err != nil {
				return nil
			}
		}
		return nil
	}
	if r.Cmd == CmdUDP {
		caddr, err := r.UDP(c, s.ServerAddr)
		if err != nil {
			return err
		}
		_, p, err := net.SplitHostPort(caddr.String())
		if err != nil {
			return err
		}
		if p == "0" {
			time.Sleep(time.Duration(s.UDPSessionTime) * time.Second)
			return nil
		}
		ch := make(chan byte)
		s.TCPUDPAssociate.Set(caddr.String(), ch, cache.DefaultExpiration)
		<-ch
		return nil
	}
	return ErrUnsupportCmd
}

// UDPHandle auto handle packet. You may prefer to do yourself.
func (h *DefaultHandle) UDPHandle(s *Server, addr *net.UDPAddr, d *Datagram) error {
	send := func(ue *UDPExchange, data []byte) error {
		_, err := ue.RemoteConn.Write(data)
		if err != nil {
			return err
		}
		if Debug {
			log.Printf("Sent UDP data to remote. client: %#v server: %#v remote: %#v data: %#v\n", ue.ClientAddr.String(), ue.RemoteConn.LocalAddr().String(), ue.RemoteConn.RemoteAddr().String(), data)
		}
		return nil
	}

	var ue *UDPExchange
	iue, ok := s.UDPExchanges.Get(addr.String())
	if ok {
		ue = iue.(*UDPExchange)
		return send(ue, d.Data)
	}

	if Debug {
		log.Printf("Call udp: %#v\n", d.Address())
	}
	c, err := Dial.Dial("udp", d.Address())
	if err != nil {
		v, ok := s.TCPUDPAssociate.Get(addr.String())
		if ok {
			ch := v.(chan byte)
			ch <- 0x00
			s.TCPUDPAssociate.Delete(addr.String())
		}
		return err
	}
	// A UDP association terminates when the TCP connection that the UDP
	// ASSOCIATE request arrived on terminates.
	rc := c.(*net.UDPConn)
	ue = &UDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
	}
	if Debug {
		log.Printf("Created remote UDP conn for client. client: %#v server: %#v remote: %#v\n", addr.String(), ue.RemoteConn.LocalAddr().String(), d.Address())
	}
	if err := send(ue, d.Data); err != nil {
		v, ok := s.TCPUDPAssociate.Get(ue.ClientAddr.String())
		if ok {
			ch := v.(chan byte)
			ch <- 0x00
			s.TCPUDPAssociate.Delete(ue.ClientAddr.String())
		}
		ue.RemoteConn.Close()
		return err
	}
	s.UDPExchanges.Set(ue.ClientAddr.String(), ue, cache.DefaultExpiration)
	go func(ue *UDPExchange) {
		defer func() {
			v, ok := s.TCPUDPAssociate.Get(ue.ClientAddr.String())
			if ok {
				ch := v.(chan byte)
				ch <- 0x00
				s.TCPUDPAssociate.Delete(ue.ClientAddr.String())
			}
			s.UDPExchanges.Delete(ue.ClientAddr.String())
			ue.RemoteConn.Close()
		}()
		var b [65536]byte
		for {
			if s.UDPDeadline != 0 {
				if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
					log.Println(err)
					break
				}
			}
			n, err := ue.RemoteConn.Read(b[:])
			if err != nil {
				break
			}
			if Debug {
				log.Printf("Got UDP data from remote. client: %#v server: %#v remote: %#v data: %#v\n", ue.ClientAddr.String(), ue.RemoteConn.LocalAddr().String(), ue.RemoteConn.RemoteAddr().String(), b[0:n])
			}
			a, addr, port, err := ParseAddress(ue.ClientAddr.String())
			if err != nil {
				log.Println(err)
				break
			}
			d1 := NewDatagram(a, addr, port, b[0:n])
			if _, err := s.UDPConn.WriteToUDP(d1.Bytes(), ue.ClientAddr); err != nil {
				break
			}
			if Debug {
				log.Printf("Sent Datagram. client: %#v server: %#v remote: %#v data: %#v %#v %#v %#v %#v %#v datagram address: %#v\n", ue.ClientAddr.String(), ue.RemoteConn.LocalAddr().String(), ue.RemoteConn.RemoteAddr().String(), d1.Rsv, d1.Frag, d1.Atyp, d1.DstAddr, d1.DstPort, d1.Data, d1.Address())
			}
		}
	}(ue)
	return nil
}
