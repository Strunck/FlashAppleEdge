package main

import (
	"Strunck/FlashApple/device"
	"fmt"
)

const Version = "dev - 0.0.2"

func main() {
	fmt.Println("FlashApple Edge Client Version ", Version)

	device.Run()

}
