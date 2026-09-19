package main

import (
	"Strunck/FlashApple/device"
	"Strunck/FlashApple/web"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var Version = "development"

func main() {
	fmt.Println("FlashApple Edge Client Version", Version)

	ctx := context.Background()
	s := &device.State{}
	errCh := make(chan error, 1)
	initDoneCh := make(chan bool)

	go device.Run(ctx, s, errCh, initDoneCh)

	ctx, userstop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer userstop()

	for {
		select {
		case <-make(chan struct{}): // This will block forever, effectively keeping the program running

		case <-initDoneCh:
			fmt.Println("Initialization done")
			go web.Run(s)

		case <-ctx.Done():
			fmt.Println("Done()", context.Cause(ctx))
			_ = s.CloseAllFiles()
			os.Exit(0)
		case err := <-errCh:
			fmt.Printf("Error main: %v", err)
		}
	}
}
