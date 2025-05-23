package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ErrCouldNotConnectToServer error = errors.New("could not connect to server")

const (
	signalCTRL_D = iota
	signalServerInterrupt
)

var verboseFlag bool
var timeoutFlag time.Duration

func init() {
	flag.DurationVar(&timeoutFlag, "timeout", time.Second*10, "timeout in seconds")
	flag.BoolVar(&verboseFlag, "v", false, "v")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	args := flag.Args()

	if len(args) < 2 {
		log.Fatal("usage: go run main.go <host> <port>")
	}

	tClient := NewTelnetClient(args[0]+":"+args[1], timeoutFlag, os.Stdin, os.Stdout)

	err := tClient.Connect()
	if err != nil {
		fmt.Println(ErrCouldNotConnectToServer)
		return
	}

	customSignal := make(chan int, 1)
	serverNotRespondSignal := make(chan int, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				select {
				case <-serverNotRespondSignal:
					log.Fatal("server not responding")
					return
				default:
					err = tClient.Send()
					if err != nil {
						if err == io.EOF {
							customSignal <- signalCTRL_D
							return
						}
						log.Fatal(err)
						continue
					}
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
						serverNotRespondSignal <- signalServerInterrupt
						return
					}
					log.Fatal(err)
					continue
				}
			}
		}
	}()

	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT)
	select {
	case s := <-customSignal:
		if s == signalCTRL_D {
			log.Fatal("program exited with CTRL+D")
		}
	case <-ctx.Done():
		log.Fatal("program exited CTRL+C")
	}
	cancel()
	tClient.Close()
}
