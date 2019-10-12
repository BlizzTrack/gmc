package main

import (
	"github.com/blizztrack/gmc"
	"log"
)

// TODO: have config for user control
func main() {
	if err := gmc.NewServer(":11211", 10000); err != nil {
		log.Panic(err)
	}
}
