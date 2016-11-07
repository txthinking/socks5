package socks5

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"net"
	"time"
)

const AES_256_KEY_LENTH = 32

// Cut or append empty data on the key.
// make the key length equal 32
func makeAES256Key(k []byte) []byte {
	if len(k) < AES_256_KEY_LENTH {
		var a []byte = make([]byte, AES_256_KEY_LENTH-len(k))
		return append(k, a...)
	}
	if len(k) > AES_256_KEY_LENTH {
		return k[:AES_256_KEY_LENTH]
	}
	return k
}

type CipherReadWriter struct {
	c  net.Conn
	sr cipher.StreamReader
	sw cipher.StreamWriter
}

func NewCipherReadWriter(c net.Conn, key []byte, iv []byte) (*CipherReadWriter, error) {
	if len(iv) != aes.BlockSize {
		return nil, errors.New("Invalid IV length")
	}
	block, err := aes.NewCipher(makeAES256Key(key))
	if err != nil {
		return nil, err
	}
	return &CipherReadWriter{
		c: c,
		sr: cipher.StreamReader{
			S: cipher.NewCFBDecrypter(block, iv),
			R: c,
		},
		sw: cipher.StreamWriter{
			S: cipher.NewCFBEncrypter(block, iv),
			W: c,
		},
	}, nil
}

func (c *CipherReadWriter) Read(b []byte) (n int, err error) {
	return c.sr.Read(b)
}

func (c *CipherReadWriter) Write(b []byte) (n int, err error) {
	return c.sw.Write(b)
}

func (c *CipherReadWriter) Close() error {
	return c.c.Close()
}

func (c *CipherReadWriter) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}
func (c *CipherReadWriter) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}
func (c *CipherReadWriter) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}
func (c *CipherReadWriter) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}
func (c *CipherReadWriter) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}
