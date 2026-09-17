package main

import (
	"Strunck/FlashApple/device"
	"fmt"
)

var Version = "development"

func main() {
	fmt.Println("FlashApple Edge Client Version", Version)

	device.Run()

}
