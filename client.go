package socks5

import (
	"errors"
	"net"
	"time"
)

// Client is socks5 client wrapper
type Client struct {
	Server        string
	UserName      string
	Password      string
	TCPConn       *net.TCPConn
	UDPConn       *net.UDPConn
	RemoteAddress net.Addr
	TCPDeadline   int
	TCPTimeout    int
	UDPDeadline   int
}

// This is just create a client, you need to use Dial to create conn
func NewClient(addr, username, password string, tcpTimeout, tcpDeadline, udpDeadline int) (*Client, error) {
	c := &Client{
		Server:      addr,
		UserName:    username,
		Password:    password,
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
	}
	return c, nil
}

func (c *Client) Dial(network, addr string) (net.Conn, error) {
	c = &Client{
		Server:      c.Server,
		UserName:    c.UserName,
		Password:    c.Password,
		TCPTimeout:  c.TCPTimeout,
		TCPDeadline: c.TCPDeadline,
		UDPDeadline: c.UDPDeadline,
	}
	if network == "tcp" {
		var err error
		c.RemoteAddress, err = net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, err
		}
		if err := c.Negotiate(); err != nil {
			return nil, err
		}
		a, h, p, err := ParseAddress(addr)
		if err != nil {
			return nil, err
		}
		if a == ATYPDomain {
			h = h[1:]
		}
		if _, err := c.Request(NewRequest(CmdConnect, a, h, p)); err != nil {
			return nil, err
		}
		return c, nil
	}
	if network == "udp" {
		var err error
		c.RemoteAddress, err = net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return nil, err
		}
		if err := c.Negotiate(); err != nil {
			return nil, err
		}

		// TODO support local udp addr
		a, h, p, err := ParseAddress(addr)
		if err != nil {
			return nil, err
		}
		if a == ATYPIPv4 || a == ATYPDomain {
			a = ATYPIPv4
			h = net.IPv4zero
		}
		if a == ATYPIPv6 {
			h = net.IPv6zero
		}
		p = []byte{0x00, 0x00}
		rp, err := c.Request(NewRequest(CmdUDP, a, h, p))
		if err != nil {
			return nil, err
		}
		tmp, err := Dial.Dial("udp", rp.Address())
		if err != nil {
			return nil, err
		}
		c.UDPConn = tmp.(*net.UDPConn)
		return c, nil
	}
	return nil, errors.New("unsupport network")
}

func (c *Client) Read(b []byte) (int, error) {
	if c.UDPConn == nil {
		return c.TCPConn.Read(b)
	}
	b1 := make([]byte, 65535)
	n, err := c.UDPConn.Read(b1)
	if err != nil {
		return 0, err
	}
	d, err := NewDatagramFromBytes(b1[0:n])
	if err != nil {
		return 0, err
	}
	if len(b) < len(d.Data) {
		return 0, errors.New("b too small")
	}
	n = copy(b, d.Data)
	return n, nil
}

func (c *Client) Write(b []byte) (int, error) {
	if c.UDPConn == nil {
		return c.TCPConn.Write(b)
	}
	a, h, p, err := ParseAddress(c.RemoteAddress.String())
	if err != nil {
		return 0, err
	}
	if a == ATYPDomain {
		h = h[1:]
	}
	d := NewDatagram(a, h, p, b)
	b1 := d.Bytes()
	n, err := c.UDPConn.Write(b1)
	if err != nil {
		return 0, err
	}
	if len(b1) != n {
		return 0, errors.New("not write full")
	}
	return len(b), nil
}

func (c *Client) Close() error {
	if c.UDPConn == nil {
		return c.TCPConn.Close()
	}
	c.TCPConn.Close()
	return c.UDPConn.Close()
}

func (c *Client) LocalAddr() net.Addr {
	if c.UDPConn == nil {
		return c.TCPConn.LocalAddr()
	}
	return c.UDPConn.LocalAddr()
}

func (c *Client) RemoteAddr() net.Addr {
	return c.RemoteAddress
}

func (c *Client) SetDeadline(t time.Time) error {
	if c.UDPConn == nil {
		return c.TCPConn.SetDeadline(t)
	}
	return c.UDPConn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	if c.UDPConn == nil {
		return c.TCPConn.SetReadDeadline(t)
	}
	return c.UDPConn.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	if c.UDPConn == nil {
		return c.TCPConn.SetWriteDeadline(t)
	}
	return c.UDPConn.SetWriteDeadline(t)
}

func (c *Client) Negotiate() error {
	con, err := Dial.Dial("tcp", c.Server)
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
