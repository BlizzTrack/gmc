package main

import (
	"github.com/blizztrack/gmc"
	"log"
)

func main() {
	if err := gmc.NewServer(":11211"); err != nil {
		log.Panic(err)
	}
}
