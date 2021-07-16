## socks5

[English](README.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/socks5)](https://goreportcard.com/report/github.com/txthinking/socks5)
[![GoDoc](https://godoc.org/github.com/txthinking/socks5?status.svg)](https://godoc.org/github.com/txthinking/socks5)
[![捐赠](https://img.shields.io/badge/%E6%94%AF%E6%8C%81-%E6%8D%90%E8%B5%A0-ff69b4.svg)](https://github.com/sponsors/txthinking)
[![交流群](https://img.shields.io/badge/%E7%94%B3%E8%AF%B7%E5%8A%A0%E5%85%A5-%E4%BA%A4%E6%B5%81%E7%BE%A4-ff69b4.svg)](https://docs.google.com/forms/d/e/1FAIpQLSdzMwPtDue3QoezXSKfhW88BXp57wkbDXnLaqokJqLeSWP9vQ/viewform)

SOCKS Protocol Version 5 Library.

完整 TCP/UDP 和 IPv4/IPv6 支持.
目标: KISS, less is more, small API, code is like the original protocol.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### 获取
```
$ go get github.com/txthinking/socks5
```

### Struct的概念 对标 原始协议里的概念

* Negotiation:
    * `type NegotiationRequest struct`
        * `func NewNegotiationRequest(methods []byte)`, in client
        * `func (r *NegotiationRequest) WriteTo(w io.Writer)`, client writes to server
        * `func NewNegotiationRequestFrom(r io.Reader)`, server reads from client
    * `type NegotiationReply struct`
        * `func NewNegotiationReply(method byte)`, in server
        * `func (r *NegotiationReply) WriteTo(w io.Writer)`, server writes to client
        * `func NewNegotiationReplyFrom(r io.Reader)`, client reads from server
* User and password negotiation:
    * `type UserPassNegotiationRequest struct`
        * `func NewUserPassNegotiationRequest(username []byte, password []byte)`, in client
        * `func (r *UserPassNegotiationRequest) WriteTo(w io.Writer)`, client writes to server
        * `func NewUserPassNegotiationRequestFrom(r io.Reader)`, server reads from client
    * `type UserPassNegotiationReply struct`
        * `func NewUserPassNegotiationReply(status byte)`, in server
        * `func (r *UserPassNegotiationReply) WriteTo(w io.Writer)`, server writes to client
        * `func NewUserPassNegotiationReplyFrom(r io.Reader)`, client reads from server
* Request:
    * `type Request struct`
        * `func NewRequest(cmd byte, atyp byte, dstaddr []byte, dstport []byte)`, in client
        * `func (r *Request) WriteTo(w io.Writer)`, client writes to server
        * `func NewRequestFrom(r io.Reader)`, server reads from client
        * After server gets the client's *Request, processes...
* Reply:
    * `type Reply struct`
        * `func NewReply(rep byte, atyp byte, bndaddr []byte, bndport []byte)`, in server
        * `func (r *Reply) WriteTo(w io.Writer)`, server writes to client
        * `func NewReplyFrom(r io.Reader)`, client reads from server
* Datagram:
    * `type Datagram struct`
        * `func NewDatagram(atyp byte, dstaddr []byte, dstport []byte, data []byte)`
        * `func NewDatagramFromBytes(bb []byte)`
        * `func (d *Datagram) Bytes()`

### 高级 API

**Server**. 你可以自己处理client请求在读取**Request**后. 同时, 这里有一个高级接口

* `type Server struct`
* `type Handler interface`
    * `TCPHandle(*Server, *net.TCPConn, *Request) error`
    * `UDPHandle(*Server, *net.UDPAddr, *Datagram) error`

举例:

```
s, _ := NewClassicServer(addr, ip, username, password, tcpTimeout, udpTimeout)
s.ListenAndServe(Handler)
```

* 如果你想要一个标准socks5 server, 传入nil即可
* 如果你想要自己处理请求, 传入一个你自己的Handler

**Client**. 这里有个socks5 client, 支持TCP和UDP, 返回net.Conn.

* `type Client struct`

举例:

```
c, _ := socks5.NewClient(server, username, password, tcpTimeout, udpTimeout)
conn, _ := c.Dial(network, addr)
```

### 用户:

 * Brook [https://github.com/txthinking/brook](https://github.com/txthinking/brook)
 * Shiliew [https://www.shiliew.com](https://www.shiliew.com)

## 开源协议

基于 MIT 协议开源
