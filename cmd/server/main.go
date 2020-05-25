package main

import (
	"log"

	"github.com/seanpfeifer/hostrelay/reliable"
	"github.com/seanpfeifer/hostrelay/unreliable"
)

const (
	defaultTCPHost = ":8080"
	defaultUDPHost = ":8585"
)

func main() {
	log.Println("Starting server...")

	go func() {
		err := unreliable.ListenAndServeUDP("udp", defaultUDPHost)
		FatalOnError(err)
	}()

	err := reliable.ListenAndServeTCP("tcp", defaultTCPHost)
	FatalOnError(err)
}
