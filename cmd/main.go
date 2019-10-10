package main

import (
	"github.com/blizztrack/gmc"
	"log"
)

func main() {
	if err := gmc.NewServer(":11212"); err != nil {
		log.Panic(err)
	}
}
