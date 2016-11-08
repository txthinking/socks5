package socks5

const (
	VER byte = 0x05

	METHOD_NONE              byte = 0x00
	METHOD_GSSAPI            byte = 0x01 // MUST support // todo
	METHOD_USERNAME_PASSWORD byte = 0x02 // SHOULD support
	METHOD_UNSUPPORT_ALL     byte = 0xFF

	USER_PASS_VER            byte = 0x01
	USER_PASS_STATUS_SUCCESS byte = 0x00
	USER_PASS_STATUS_FAILURE byte = 0x01 // just other than 0x00

	CMD_CONNECT byte = 0x01
	CMD_BIND    byte = 0x02
	CMD_UDP     byte = 0x03

	ATYP_IPV4   byte = 0x01 // 4 octets
	ATYP_DOMAIN byte = 0x03 // The first octet of the address field contains the number of octets of name that follow, there is no terminating NUL octet.
	ATYP_IPV6   byte = 0x04 // 16 octets

	REP_SUCCESS               byte = 0x00
	REP_SERVER_FAILURE        byte = 0x01
	REP_NOT_ALLOWED           byte = 0x02
	REP_NETWORK_UNREACHABLE   byte = 0x03
	REP_HOST_UNREACHABLE      byte = 0x04
	REP_CONNECTION_REFUSED    byte = 0x05
	REP_TTL_EXPIRED           byte = 0x06
	REP_COMMAND_NOT_SUPPORTED byte = 0x07
	REP_ADDRESS_NOT_SUPPORTED byte = 0x08
)

type NegotiationRequest struct {
	Ver      byte
	NMethods byte
	Methods  []byte // 1-255 bytes
}

type NegotiationReply struct {
	Ver    byte
	Method byte
}

type UserPassNegotiationRequest struct {
	Ver    byte
	Ulen   byte
	Uname  []byte // 1-255 bytes
	Plen   byte
	Passwd []byte // 1-255 bytes
}

type UserPassNegotiationReply struct {
	Ver    byte
	Status byte
}

type Request struct {
	Ver     byte
	Cmd     byte
	Rsv     byte // 0x00
	Atyp    byte
	DstAddr []byte
	DstPort []byte // 2 bytes
}

type Reply struct {
	Ver  byte
	Rep  byte
	Rsv  byte // 0x00
	Atyp byte
	// CONNECT socks server's address which used to connect to dst addr
	// BIND ...
	// UDP socks server's address which used to connect to dst addr
	BndAddr []byte
	// CONNECT socks server's port which used to connect to dst addr
	// BIND ...
	// UDP socks server's port which used to connect to dst addr
	BndPort []byte // 2 bytes
}

type UDPHeader struct {
	Rsv     []byte // 0x00 0x00
	Flag    byte
	Atyp    byte
	DstAddr []byte
	DstPort []byte // 2 bytes
}
