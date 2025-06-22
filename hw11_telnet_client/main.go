package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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
					log.Fatal(err)
				}
				close(EOFSignal)
				return
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
					log.Fatal(err)
				}
				close(serverNotRespondSignal)
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Print("program exited CTRL+C")
		break
	case <-EOFSignal:
		log.Print("program exited with CTRL+D")
		break
	case <-serverNotRespondSignal:
		log.Print("server not respond")
		break
	}
	cancel()
	tClient.Close()
}
