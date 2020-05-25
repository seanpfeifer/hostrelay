package unreliable

import (
	"fmt"
	"log"
	"net"
)

func ListenAndServeUDP(network, address string) error {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP(network, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf(`Listening for %s on "%s"`, conn.LocalAddr().Network(), conn.LocalAddr().String())

	buf := make([]byte, 2048)

	count := 0

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("Received: %s\n", string(buf[:n]))

		conn.WriteToUDP([]byte(fmt.Sprintf("Thanks for message #%d", count)), addr)
		count++
	}
}
