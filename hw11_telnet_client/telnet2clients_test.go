package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetVS2Client(t *testing.T) {
	t.Run("2 concurrent clients against server", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		// Do NOT require no error on Close here; just defer close ignoring error
		defer func() { _ = l.Close() }()

		var serverWg sync.WaitGroup
		done := make(chan struct{})

		// Server goroutine: accepts connections until signaled to stop
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					conn, err := l.Accept()
					if err != nil {
						// Exit gracefully if listener closed
						if errors.Is(err, net.ErrClosed) {
							return
						}
						t.Logf("Accept error: %v", err)
						return
					}

					serverWg.Add(1)
					go func(c net.Conn) {
						defer serverWg.Done()
						defer c.Close()

						request := make([]byte, 1024)
						n, err := c.Read(request)
						require.NoError(t, err)
						require.Equal(t, "hello from server\n", string(request[:n]))

						n, err = c.Write([]byte("world\n"))
						require.NoError(t, err)
						require.NotEqual(t, 0, n)
					}(conn)
				}
			}
		}()

		// Prepare client goroutines
		var clientWg sync.WaitGroup
		clients := []int{0, 1} // Two clients

		for range clients {
			clientWg.Add(1)
			go func() {
				defer clientWg.Done()

				in := &bytes.Buffer{}
				out := &bytes.Buffer{}

				timeout, err := time.ParseDuration("10s")
				require.NoError(t, err)

				client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
				require.NoError(t, client.Connect())
				defer func() { require.NoError(t, client.Close()) }()

				in.WriteString("hello from server\n")
				err = client.Send()
				require.NoError(t, err)

				err = client.Receive()
				require.NoError(t, err)
				require.Equal(t, "world\n", out.String())
			}()
		}

		// Wait for all clients to finish
		clientWg.Wait()

		// Signal server to stop accepting new connections
		close(done)

		// Close the listener to unblock Accept()
		_ = l.Close()

		// Wait for all server handlers to finish
		serverWg.Wait()
	})
}
