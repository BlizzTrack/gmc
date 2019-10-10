package main

import (
	"github.com/blizztrack/gmc"
	"log"
)

func main() {
	if err := gmc.NewServer(); err != nil {
		log.Panic(err)
	}
}
