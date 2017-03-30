package socks5

import "testing"

func TestParseAddress(t *testing.T) {
	t.Log(ParseAddress("127.0.0.1:80"))
	t.Log(ParseAddress("[::1]:80"))
	t.Log(ParseAddress("a.com:80"))
}
