package socks5

import (
    "encoding/binary"
    "strconv"
    "bytes"
)

func (r *Request) Address() string {
    var addr []byte
    if r.Atyp == ATYP_DOMAIN {
        addr = r.DstAddr[1:]
    }else{
        addr = r.DstAddr
    }
    return bytes.NewBuffer(addr).String() + ":" + strconv.Itoa(int(binary.BigEndian.Uint16(r.DstPort)))
}
