package socks5

import (
	"errors"
	"net"
	"time"
)

// Client is socks5 client wrapper
type Client struct {
	Server   string
	UserName string
	Password string
	// On cmd UDP, let server control the tcp and udp connection relationship
	Conn          net.Conn
	PacketConn    net.PacketConn
	RemoteAddress net.Addr
	TCPTimeout    int
	UDPTimeout    int
}

// This is just create a client, you need to use Dial to create conn
func NewClient(addr, username, password string, tcpTimeout, udpTimeout int) (*Client, error) {
	c := &Client{
		Server:     addr,
		UserName:   username,
		Password:   password,
		TCPTimeout: tcpTimeout,
		UDPTimeout: udpTimeout,
	}
	return c, nil
}

func (c *Client) Dial(network, addr string) (net.Conn, error) {
	return c.DialWithLocalAddr(network, "", addr, nil)
}

func (c *Client) DialWithLocalAddr(network, src, dst string, remoteAddr net.Addr) (net.Conn, error) {
	c = &Client{
		Server:        c.Server,
		UserName:      c.UserName,
		Password:      c.Password,
		TCPTimeout:    c.TCPTimeout,
		UDPTimeout:    c.UDPTimeout,
		RemoteAddress: remoteAddr,
	}
	var err error
	if network == "tcp" {
		if c.RemoteAddress == nil {
			c.RemoteAddress, err = net.ResolveTCPAddr("tcp", dst)
			if err != nil {
				return nil, err
			}
		}
		var la *net.TCPAddr
		if src != "" {
			la, err = net.ResolveTCPAddr("tcp", src)
			if err != nil {
				return nil, err
			}
		}
		if err := c.Negotiate(la); err != nil {
			return nil, err
		}
		a, h, p, err := ParseAddress(dst)
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
		if c.RemoteAddress == nil {
			c.RemoteAddress, err = net.ResolveUDPAddr("udp", dst)
			if err != nil {
				return nil, err
			}
		}
		var la *net.TCPAddr
		if src != "" {
			la, err = net.ResolveTCPAddr("tcp", src)
			if err != nil {
				return nil, err
			}
		}
		if err := c.Negotiate(la); err != nil {
			return nil, err
		}

		var laddr *net.UDPAddr
		if src != "" {
			laddr, err = net.ResolveUDPAddr("udp", src)
			if err != nil {
				return nil, err
			}
		}
		if src == "" {
			laddr = &net.UDPAddr{
				IP:   c.Conn.LocalAddr().(*net.TCPAddr).IP,
				Port: c.Conn.LocalAddr().(*net.TCPAddr).Port,
				Zone: c.Conn.LocalAddr().(*net.TCPAddr).Zone,
			}
		}
		a, h, p, err := ParseAddress(laddr.String())
		if err != nil {
			return nil, err
		}
		rp, err := c.Request(NewRequest(CmdUDP, a, h, p))
		if err != nil {
			return nil, err
		}
		raddr, err := net.ResolveUDPAddr("udp", rp.Address())
		if err != nil {
			return nil, err
		}
		c.PacketConn, err = Dial.DialUDP("udp", laddr, raddr)
		if err != nil {
			return nil, err
		}
		if c.UDPTimeout != 0 {
			if err := c.PacketConn.SetDeadline(time.Now().Add(time.Duration(c.UDPTimeout) * time.Second)); err != nil {
				return nil, err
			}
		}
		return c, nil
	}
	return nil, errors.New("unsupport network")
}

func (c *Client) Read(b []byte) (int, error) {
	if c.PacketConn == nil {
		return c.Conn.Read(b)
	}
	n, _, err := c.PacketConn.ReadFrom(b)
	if err != nil {
		return 0, err
	}
	d, err := NewDatagramFromBytes(b[0:n])
	if err != nil {
		return 0, err
	}
	n = copy(b, d.Data)
	return n, nil
}

func (c *Client) Write(b []byte) (int, error) {
	if c.PacketConn == nil {
		return c.Conn.Write(b)
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
	n, err := c.PacketConn.WriteTo(b1, c.RemoteAddress)
	if err != nil {
		return 0, err
	}
	if len(b1) != n {
		return 0, errors.New("not write full")
	}
	return len(b), nil
}

func (c *Client) Close() error {
	if c.PacketConn == nil {
		return c.Conn.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
	return c.PacketConn.Close()
}

func (c *Client) LocalAddr() net.Addr {
	if c.PacketConn == nil {
		return c.Conn.LocalAddr()
	}
	return c.PacketConn.LocalAddr()
}

func (c *Client) RemoteAddr() net.Addr {
	return c.RemoteAddress
}

func (c *Client) SetDeadline(t time.Time) error {
	if c.PacketConn == nil {
		return c.Conn.SetDeadline(t)
	}
	return c.PacketConn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	if c.PacketConn == nil {
		return c.Conn.SetReadDeadline(t)
	}
	return c.PacketConn.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	if c.PacketConn == nil {
		return c.Conn.SetWriteDeadline(t)
	}
	return c.PacketConn.SetWriteDeadline(t)
}

func (c *Client) Negotiate(laddr *net.TCPAddr) error {
	raddr, err := net.ResolveTCPAddr("tcp", c.Server)
	if err != nil {
		return err
	}
	c.Conn, err = Dial.DialTCP("tcp", laddr, raddr)
	if err != nil {
		return err
	}
	if c.TCPTimeout != 0 {
		if err := c.Conn.SetDeadline(time.Now().Add(time.Duration(c.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	m := MethodNone
	if c.UserName != "" && c.Password != "" {
		m = MethodUsernamePassword
	}
	rq := NewNegotiationRequest([]byte{m})
	if _, err := rq.WriteTo(c.Conn); err != nil {
		return err
	}
	rp, err := NewNegotiationReplyFrom(c.Conn)
	if err != nil {
		return err
	}
	if rp.Method != m {
		return errors.New("Unsupport method")
	}
	if m == MethodUsernamePassword {
		urq := NewUserPassNegotiationRequest([]byte(c.UserName), []byte(c.Password))
		if _, err := urq.WriteTo(c.Conn); err != nil {
			return err
		}
		urp, err := NewUserPassNegotiationReplyFrom(c.Conn)
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
	if _, err := r.WriteTo(c.Conn); err != nil {
		return nil, err
	}
	rp, err := NewReplyFrom(c.Conn)
	if err != nil {
		return nil, err
	}
	if rp.Rep != RepSuccess {
		return nil, errors.New("Host unreachable")
	}
	return rp, nil
}
