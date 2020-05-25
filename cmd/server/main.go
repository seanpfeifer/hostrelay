package main

import (
	"log"

	"github.com/seanpfeifer/hostrelay/server"
)

const (
	defaultTCPHost = ":8080"
	defaultUDPHost = ":8585"
)

func main() {
	log.Println("Starting server...")

	go func() {
		//err := server.ListenAndServe("udp", defaultUDPHost)
		//FatalOnError(err)
	}()

	err := server.ListenAndServeTCP("tcp", defaultTCPHost)
	FatalOnError(err)
}
