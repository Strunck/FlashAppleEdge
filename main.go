package main

import (
	"Strunck/FlashApple/device"
	"fmt"
)

const Version = "0.0.1"

func main() {
	fmt.Println("FlashApple Edge Client Version ", Version)

	device.Run()

}
