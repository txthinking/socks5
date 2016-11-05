package socks5

import (
    "errors"
    "io"
    "log"
)

var (
    ERROR_BAD_REPLY = errors.New("Bad Reply")
)

type Client struct {
}

func NewNegotiationRequest(methods []byte) *NegotiationRequest {
    return &NegotiationRequest{
        Ver:      VER,
        NMethods: byte(len(methods)),
        Methods:  methods,
    }
}

func (r *NegotiationRequest) Write(w io.Writer) error {
    if _, err := w.Write([]byte{r.Ver}); err != nil {
        return err
    }
    if _, err := w.Write([]byte{r.NMethods}); err != nil {
        return err
    }
    if _, err := w.Write(r.Methods); err != nil {
        return err
    }
    log.Printf("Sent NegotiationRequest: %#v %#v %#v\n", r.Ver, r.NMethods, r.Methods)
    return nil
}

func NewNegotiationReplyFrom(r io.Reader) (*NegotiationReply, error) {
    bb := make([]byte, 2)
    if _, err := io.ReadFull(r, bb); err != nil {
        return nil, err
    }
    if bb[0] != VER {
        return nil, ERROR_VERSION
    }
    log.Printf("Got NegotiationReply: %#v %#v\n", bb[0], bb[1])
    return &NegotiationReply{
        Ver:    bb[0],
        Method: bb[1],
    }, nil
}

func NewUserPassNegotiationRequest(username []byte, password []byte) *UserPassNegotiationRequest {
    return &UserPassNegotiationRequest{
        Ver:    USER_PASS_VER,
        Ulen:   byte(len(username)),
        Uname:  username,
        Plen:   byte(len(password)),
        Passwd: password,
    }
}

func (r *UserPassNegotiationRequest) WriteTo(w io.Writer) error {
    if _, err := w.Write([]byte{r.Ver, r.Ulen}); err != nil {
        return err
    }
    if _, err := w.Write(r.Uname); err != nil {
        return err
    }
    if _, err := w.Write([]byte{r.Plen}); err != nil {
        return err
    }
    if _, err := w.Write(r.Passwd); err != nil {
        return err
    }
    log.Printf("Sent UserNameNegotiationRequest: %#v %#v %#v %#v %#v\n", r.Ver, r.Ulen, r.Uname, r.Plen, r.Passwd)
    return nil
}

func NewUserPassNegotiationReplyFrom(r io.Reader) (*UserPassNegotiationReply, error) {
    bb := make([]byte, 2)
    if _, err := io.ReadFull(r, bb); err != nil {
        return nil, err
    }
    if bb[0] != USER_PASS_VER {
        return nil, ERROR_USER_PASS_VERSION
    }
    log.Printf("Got UserPassNegotiationReply: %#v %#v \n", bb[0], bb[1])
    return &UserPassNegotiationReply{
        Ver:    bb[0],
        Status: bb[1],
    }, nil
}

func NewRequest(cmd byte, atyp byte, dstaddr []byte, dstport []byte) *Request {
    if atyp == ATYP_DOMAIN {
        dstaddr = append([]byte{byte(len(dstaddr))}, dstaddr...)
    }
    return &Request{
        Ver:     VER,
        Cmd:     cmd,
        Rsv:     0x00,
        Atyp:    atyp,
        DstAddr: dstaddr,
        DstPort: dstport,
    }
}

func (r *Request) WriteTo(w io.Writer) error {
    if _, err := w.Write([]byte{r.Ver, r.Cmd, r.Rsv, r.Atyp}); err != nil {
        return err
    }
    if r.Atyp == ATYP_DOMAIN {
        if _, err := w.Write([]byte{byte(len(r.DstAddr))}); err != nil {
            return err
        }
    }
    if _, err := w.Write(r.DstAddr); err != nil {
        return err
    }
    if _, err := w.Write(r.DstPort); err != nil {
        return err
    }
    log.Printf("Sent Request: %#v %#v %#v %#v %#v %#v\n", r.Ver, r.Cmd, r.Rsv, r.Atyp, r.DstAddr, r.DstPort)
    return nil
}

func NewReplyFrom(r io.Reader) (*Reply, error) {
    bb := make([]byte, 4)
    if _, err := io.ReadFull(r, bb); err != nil {
        return nil, err
    }
    if bb[0] != VER {
        return nil, ERROR_VERSION
    }
    var addr []byte
    if bb[3] == ATYP_IPV4 {
        addr = make([]byte, 4)
        if _, err := io.ReadFull(r, addr); err != nil {
            return nil, err
        }
    } else if bb[3] == ATYP_IPV6 {
        addr = make([]byte, 16)
        if _, err := io.ReadFull(r, addr); err != nil {
            return nil, err
        }
    } else if bb[3] == ATYP_DOMAIN {
        dal := make([]byte, 1)
        if _, err := io.ReadFull(r, dal); err != nil {
            return nil, err
        }
        if dal[0] == 0 {
            return nil, ERROR_BAD_REPLY
        }
        addr = make([]byte, int(dal[0]))
        if _, err := io.ReadFull(r, addr); err != nil {
            return nil, err
        }
        addr = append(dal, addr...)
    } else {
        return nil, ERROR_BAD_REPLY
    }
    port := make([]byte, 2)
    if _, err := io.ReadFull(r, port); err != nil {
        return nil, err
    }
    log.Printf("Got Reply: %#v %#v %#v %#v %#v %#v\n", bb[0], bb[1], bb[2], bb[3], addr, port)
    return &Reply{
        Ver:     bb[0],
        Rep:     bb[1],
        Rsv:     bb[2],
        Atyp:    bb[3],
        BndAddr: addr,
        BndPort: port,
    }, nil
}
