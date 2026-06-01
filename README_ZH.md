## socks5

[English](README.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/socks5)](https://goreportcard.com/report/github.com/txthinking/socks5)
[![GoDoc](https://godoc.org/github.com/txthinking/socks5?status.svg)](https://godoc.org/github.com/txthinking/socks5)

[🗣 News](https://t.me/s/txthinking_news)
[🩸 Youtube](https://www.youtube.com/txthinking)

SOCKS Protocol Version 5 Library.

完整 TCP/UDP 和 IPv4/IPv6 支持.
目标: KISS, less is more, small API, code is like the original protocol.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### 获取
```
$ go get github.com/txthinking/txthinkingsocks5
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

> 这可以满足经典场景，特殊场景推荐你选择上面的小API来自定义。

**Server**: 支持UDP和TCP

* `type Server struct`
* `type Handler interface`
    * `TCPHandle(*Server, *net.TCPConn, *Request) error`
    * `UDPHandle(*Server, *net.UDPAddr, *Datagram) error`

举例:

```
server, _ := NewClassicServer(addr, ip, username, password, tcpTimeout, udpTimeout)
server.ListenAndServe(Handler)
```

**Client**: 支持TCP和UDP, 返回net.Conn

* `type Client struct`

举例:

```
client, _ := socks5.NewClient(server, username, password, tcpTimeout, udpTimeout)
conn, _ := client.Dial(network, addr)
```


### 谁在使用此项目

-   Brook: https://github.com/txthinking/brook
-   Shiliew: https://www.txthinking.com/shiliew.html
-   dismap: https://github.com/zhzyker/dismap
-   emp3r0r: https://github.com/jm33-m0/emp3r0r
-   hysteria: https://github.com/apernet/hysteria
-   mtg: https://github.com/9seconds/mtg
-   trojan-go: https://github.com/p4gefau1t/trojan-go

## 开源协议

基于 MIT 协议开源
