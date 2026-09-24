package main

import (
	"Strunck/FlashApple/device"
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

	go device.Run(ctx)

	ctx, userstop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer userstop()

	for {
		select {
		case <-make(chan struct{}): // This will block forever, effectively keeping the program running
		case <-ctx.Done():
			fmt.Println("Done()", context.Cause(ctx))

			os.Exit(0)
		}
	}
}
