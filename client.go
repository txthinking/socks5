package socks5

import (
	"errors"
	"net"
	"time"
)

// Client is socks5 client wrapper
type Client struct {
	UserName    string
	Password    string
	TCPAddr     *net.TCPAddr
	TCPConn     *net.TCPConn
	UDPAddr     *net.UDPAddr
	TCPDeadline int // not refreshed
	TCPTimeout  int
	UDPDeadline int // refreshed
}

func NewClient(addr, username, password string, tcpTimeout, tcpDeadline, udpDeadline int) (*Client, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		UserName:    username,
		Password:    password,
		TCPAddr:     taddr,
		UDPAddr:     uaddr,
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
	}
	return c, nil
}

func (c *Client) Negotiate() error {
	con, err := Dial.Dial("tcp", c.TCPAddr.String())
	if err != nil {
		return err
	}
	c.TCPConn = con.(*net.TCPConn)
	if c.TCPTimeout != 0 {
		if err := c.TCPConn.SetKeepAlivePeriod(time.Duration(c.TCPTimeout) * time.Second); err != nil {
			return err
		}
	}
	if c.TCPDeadline != 0 {
		if err := c.TCPConn.SetDeadline(time.Now().Add(time.Duration(c.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	m := MethodNone
	if c.UserName != "" && c.Password != "" {
		m = MethodUsernamePassword
	}
	rq := NewNegotiationRequest([]byte{m})
	if _, err := rq.WriteTo(c.TCPConn); err != nil {
		return err
	}
	rp, err := NewNegotiationReplyFrom(c.TCPConn)
	if err != nil {
		return err
	}
	if rp.Method != m {
		return errors.New("Unsupport method")
	}
	if m == MethodUsernamePassword {
		urq := NewUserPassNegotiationRequest([]byte(c.UserName), []byte(c.Password))
		if _, err := urq.WriteTo(c.TCPConn); err != nil {
			return err
		}
		urp, err := NewUserPassNegotiationReplyFrom(c.TCPConn)
		if err != nil {
			return err
		}
		if urp.Status != UserPassStatusSuccess {
			return ErrUserPassAuth
		}
	}
	return nil
}

func (c *Client) Request(r *Request) (*Reply, error) {
	if _, err := r.WriteTo(c.TCPConn); err != nil {
		return nil, err
	}
	rp, err := NewReplyFrom(c.TCPConn)
	if err != nil {
		return nil, err
	}
	if rp.Rep != RepSuccess {
		return nil, errors.New("Host unreachable")
	}
	return rp, nil
}
