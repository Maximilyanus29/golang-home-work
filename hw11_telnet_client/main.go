package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	ErrCouldNotConnectToServer = errors.New("could not connect to server")
	verboseFlag                bool
	timeoutFlag                time.Duration
)

func main() {
	flag.DurationVar(&timeoutFlag, "timeout", time.Second*10, "timeout in seconds")
	flag.BoolVar(&verboseFlag, "v", false, "v")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("usage: go run main.go [--timeout=10s] [-v] <host> <port>")
	}

	tClient := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeoutFlag, os.Stdin, os.Stdout)

	err := tClient.Connect()

	if err != nil {
		fmt.Println(ErrCouldNotConnectToServer)
		return
	}
	fmt.Println("connected")

	EOFSignal := make(chan byte, 1)
	serverNotRespondSignal := make(chan byte, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err = tClient.Send()
				if err != nil {
					if err == io.EOF {
						close(EOFSignal)
						return
					}
					log.Fatal(err)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err = tClient.Receive()
				if err != nil {
					if err == io.EOF {
						close(serverNotRespondSignal)
						return
					}
					log.Fatal(err)
				}
			}
		}
	}()

	select {
	case <-serverNotRespondSignal:
		fmt.Print("server not responding")
		break
	case <-ctx.Done():
		fmt.Print("program exited CTRL+C")
		break
	case <-EOFSignal:
		fmt.Print("program exited with CTRL+D")
		break
	case <-serverNotRespondSignal:
		fmt.Print("server not respond")
		break
	}
	cancel()
	tClient.Close()
}
