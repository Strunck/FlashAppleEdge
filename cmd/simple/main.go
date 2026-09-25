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

	ctx, userstop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer userstop()

	initDone := make(chan bool, 1)

	// Bootstap
	se, err := device.Bootstrap()
	if err != nil {
		fmt.Printf("Bootstrap error: %v\n", err)
		os.Exit(1)
	}
	initDone <- true

	for {
		select {
		case <-initDone:
			fmt.Println("Initialization done")
			se.Run(ctx)
		case <-make(chan struct{}): // This will block forever, effectively keeping the program running
		case <-ctx.Done():
			fmt.Println("Done()", context.Cause(ctx))
			os.Exit(0)
		}
	}
}
