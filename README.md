## socks5

[![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/socks5)](https://goreportcard.com/report/github.com/txthinking/socks5)
[![GoDoc](https://godoc.org/github.com/txthinking/socks5?status.svg)](https://godoc.org/github.com/txthinking/socks5)

SOCKS Protocol Version 5 Library.

Full TCP/UDP and IPv4/IPv6 support.
Goals: KISS, less is more, small API, code is like the original protocol.

### Install
```
$ go get github.com/txthinking/socks5
```

### Struct is like concept in protocol

* Negotiation:
    * `type NegotiationRequest struct`
        * `func NewNegotiationRequest(methods []byte) *NegotiationRequest`, in client
        * `func (r *NegotiationRequest) WriteTo(w *net.TCPConn) error`, client writes to server
        * `func NewNegotiationRequestFrom(r *net.TCPConn) (*NegotiationRequest, error)`, server reads from client
    * `type NegotiationReply struct`
        * `func NewNegotiationReply(method byte) *NegotiationReply`, in server
        * `func (r *NegotiationReply) WriteTo(w *net.TCPConn) error`, server writes to client
        * `func NewNegotiationReplyFrom(r *net.TCPConn) (*NegotiationReply, error)`, client reads from server
* User and password negotiation:
    * `type UserPassNegotiationRequest struct`
        * `func NewUserPassNegotiationRequest(username []byte, password []byte) *UserPassNegotiationRequest`, in client
        * `func (r *UserPassNegotiationRequest) WriteTo(w *net.TCPConn) error`, client writes to server
        * `func NewUserPassNegotiationRequestFrom(r *net.TCPConn) (*UserPassNegotiationRequest, error)`, server reads from client
    * `type UserPassNegotiationReply struct`
        * `func NewUserPassNegotiationReply(status byte) *UserPassNegotiationReply`, in server
        * `func (r *UserPassNegotiationReply) WriteTo(w *net.TCPConn) error`, server writes to client
        * `func NewUserPassNegotiationReplyFrom(r *net.TCPConn) (*UserPassNegotiationReply, error)`, client reads from server
* Request:
    * `type Request struct`
        * `func NewRequest(cmd byte, atyp byte, dstaddr []byte, dstport []byte) *Request`, in client
        * `func (r *Request) WriteTo(w *net.TCPConn) error`, client writes to server
        * `func NewRequestFrom(r *net.TCPConn) (*Request, error)`, server reads from client
        * After server gets the client's *Request, processes...
* Reply:
    * `type Reply struct`
        * `func NewReply(rep byte, atyp byte, bndaddr []byte, bndport []byte) *Reply`
        * `func (r *Reply) WriteTo(w *net.TCPConn) error`, server writes to client
        * `func NewReplyFrom(r *net.TCPConn) (*Reply, error)`, client reads from server
* Datagram:
    * `type Datagram struct`
        * `func NewDatagramFromBytes(bb []byte) (*Datagram, error)`, in server
        * `func NewDatagram(atyp byte, dstaddr []byte, dstport []byte, data []byte) *Datagram`, in server
        * `func (d *Datagram) Bytes() []byte`, in server

### Advanced API

You can process client's request by yourself after reading *Request from client.
Also, here is a advanced interfaces.

* `type Server struct`
* `type Handler interface`
    * `TCPHandle(*Server, *net.TCPConn, *Request) error`
    * `UDPHandle(*Server, *net.UDPAddr, *Datagram) error`

This is example:

```
s, _ := NewClassicServer(addr, ip, username, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
s.Run(Handler)
```

* If you want a standard socks5 server, pass in nil
* If you want to handle data by yourself, pass in a custom Handler


### Users:

 * Brook [https://github.com/txthinking/brook](https://github.com/txthinking/brook)
