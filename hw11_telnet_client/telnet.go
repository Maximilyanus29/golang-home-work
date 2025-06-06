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
	for scanner.Scan() {
		_, err := v.connection.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
		// Обработка строки
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil

	// EOF здесь обрабатывается автоматически
}

func (v *telnetClient) Receive() error {
	scanner := bufio.NewScanner(v.connection)
	for scanner.Scan() {
		_, err := v.out.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
		// Обработка строки
	}
	if err := scanner.Err(); err != nil {
		// Обработка ошибок (кроме EOF)
		return err
	}
	return nil
}

func ReadAllCustom(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			return b, err
		}

		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
	}
}
