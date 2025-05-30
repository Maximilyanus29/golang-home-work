package main

import (
	"bufio"
	"io"
	"log"
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
	address    string
	timeout    time.Duration
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (v *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", v.address, v.timeout)
	if err != nil {
		return err
	}
	v.connection = conn
	return nil
}

func (v *telnetClient) Close() error {
	if err := v.in.Close(); err != nil {
		log.Fatal("could not closed readerclosed")
	}
	return v.connection.Close()
}

func (v *telnetClient) Send() error {
	scanner := bufio.NewScanner(v.in)
	if scanner.Scan() {
		n, err := v.connection.Write(append(scanner.Bytes(), '\n'))
		if verboseFlag {
			log.Printf("передано %d байт\n", n)
		}
		return err
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return io.EOF
}

func (v *telnetClient) Receive() error {
	scanner := bufio.NewScanner(v.connection)
	if scanner.Scan() {
		n, err := v.out.Write(append(scanner.Bytes(), '\n'))
		if verboseFlag {
			log.Printf("получено %d байт\n", n)
		}
		return err
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return io.EOF
}
