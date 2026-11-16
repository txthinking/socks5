package socks5

import (
	"bytes"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestServerNegotiateRejectsUnsupportedMethods(t *testing.T) {
	for _, method := range []byte{MethodNone, MethodUsernamePassword} {
		t.Run(string([]byte{'0' + method}), func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			result := make(chan error, 1)
			go func() {
				conn, err := listener.Accept()
				if err != nil {
					result <- err
					return
				}
				defer conn.Close()
				conn.SetDeadline(time.Now().Add(5 * time.Second))
				server := &Server{Method: method, UserName: "user", Password: "pass"}
				result <- server.Negotiate(conn)
			}()
			conn, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Fatal(err)
			}
			request := NewNegotiationRequest([]byte{MethodGSSAPI})
			if _, err := request.WriteTo(conn); err != nil {
				t.Fatal(err)
			}
			// Keep the write half open: a rejected peer must not need to
			// send credentials or close its socket to finish negotiation.
			reply, err := io.ReadAll(conn)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(reply, []byte{Ver, MethodUnsupportAll}) {
				t.Errorf("reply = %x, want one method-rejection reply", reply)
			}
			if err := <-result; err == nil {
				t.Error("Negotiate returned nil after rejecting all offered methods")
			} else if err, ok := err.(net.Error); ok && err.Timeout() {
				t.Error("Negotiate waited for credentials after rejecting all offered methods")
			}
		})
	}
}

type negotiationStream struct {
	input  *bytes.Reader
	output bytes.Buffer
}

func (s *negotiationStream) Read(p []byte) (int, error)  { return s.input.Read(p) }
func (s *negotiationStream) Write(p []byte) (int, error) { return s.output.Write(p) }

func TestServerNegotiateSupportedMethods(t *testing.T) {
	for _, test := range []struct {
		name         string
		method       byte
		input, reply []byte
		wantError    error
	}{
		{name: "no authentication", method: MethodNone,
			input: []byte{Ver, 2, MethodGSSAPI, MethodNone}, reply: []byte{Ver, MethodNone}},
		{name: "username and password", method: MethodUsernamePassword,
			input: []byte{Ver, 1, MethodUsernamePassword, UserPassVer, 1, 'u', 1, 'p'},
			reply: []byte{Ver, MethodUsernamePassword, UserPassVer, UserPassStatusSuccess}},
		{name: "incorrect password", method: MethodUsernamePassword,
			input: []byte{Ver, 1, MethodUsernamePassword, UserPassVer, 1, 'u', 1, 'x'},
			reply: []byte{Ver, MethodUsernamePassword, UserPassVer, UserPassStatusFailure}, wantError: ErrUserPassAuth},
	} {
		t.Run(test.name, func(t *testing.T) {
			stream := &negotiationStream{input: bytes.NewReader(test.input)}
			server := &Server{Method: test.method, UserName: "u", Password: "p"}
			if err := server.Negotiate(stream); !errors.Is(err, test.wantError) {
				t.Errorf("error = %v, want %v", err, test.wantError)
			}
			if !bytes.Equal(stream.output.Bytes(), test.reply) {
				t.Errorf("reply = %x, want %x", stream.output.Bytes(), test.reply)
			}
		})
	}
}

type negotiationWriteFailure struct {
	input *bytes.Reader
	err   error
}

func (s *negotiationWriteFailure) Read(p []byte) (int, error)  { return s.input.Read(p) }
func (s *negotiationWriteFailure) Write(p []byte) (int, error) { return 0, s.err }

func TestServerNegotiatePreservesRejectionWriteError(t *testing.T) {
	want := errors.New("write failed")
	stream := &negotiationWriteFailure{input: bytes.NewReader([]byte{Ver, 1, MethodGSSAPI}), err: want}
	server := &Server{Method: MethodNone}
	if err := server.Negotiate(stream); err != want {
		t.Fatalf("error = %v, want original write error", err)
	}
}
