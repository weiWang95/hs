package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hs/cmd"

	"github.com/sirupsen/logrus"
)

func main() {
	ipv6 := flag.Bool("ipv6", false, "use ipv6")
	host := flag.String("host", "127.0.0.1", "host")
	port := flag.Int("port", 3000, "port")
	mode := flag.String("mode", "server", "mode: server|client")
	logLevel := flag.Int("logLevel", int(logrus.InfoLevel), "log level")

	flag.Parse()
	if *mode == "client" {
		fmt.Printf("%s host: %s, port: %d\n", *mode, *host, *port)
	} else {
		fmt.Printf("%s port: %d, ipv6: %v\n", *mode, *port, *ipv6)
	}

	logrus.SetLevel(logrus.Level(*logLevel))
	ctx, cancel := context.WithCancel(context.Background())

	var server *cmd.Server
	var client *cmd.Client
	done := make(chan struct{})
	if *mode == "client" {
		go func() {
			defer close(done)

			client = cmd.NewClient(*host, *port)
			if err := client.Run(ctx); err != nil {
				logrus.Error(err)
			}
		}()
	} else {
		go func() {
			defer close(done)

			server = cmd.NewServer(*port, *ipv6)
			if err := server.Start(ctx); err != nil {
				logrus.Error(err)
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case <-sigChan:
		logrus.Info("Shutdown Server ...")
		cancel()
		if server != nil {
			server.Shutdown()
		}
		if client != nil {
			client.Shutdown()
		}

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			logrus.Fatal("Force exit")
		}
	case <-done:
	}
}
