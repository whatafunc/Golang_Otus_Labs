package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", c.address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		text := scanner.Text() + "\n"
		_, err := c.conn.Write([]byte(text))
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (c *telnetClient) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		_, err := io.WriteString(c.out, scanner.Text()+"\n")
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}
