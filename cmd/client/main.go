package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/seanpfeifer/hostrelay/scan"
)

const (
	hostTCP = ":8080"
	hostUDP = ":8585"
)

type listener func(conn net.Conn, done chan struct{})

func main() {
	udpFlag := flag.Bool("udp", false, "use to enable UDP comms instead of the default TCP")
	flag.Parse()
	isUDP := *udpFlag

	network, host := "tcp", hostTCP
	if isUDP {
		network, host = "udp", hostUDP
	}

	done := make(chan struct{})
	out := make(chan string)

	conn, err := net.Dial(network, host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	list := listenForMessagesTCP
	if isUDP {
		list = listenForMessagesUDP
	}

	go list(conn, done)
	go readUserInput(out)

	for {
		select {
		case <-done:
			return
		case msg, ok := <-out:
			if !ok {
				return
			}
			send(conn, msg, isUDP)
		}
	}

}

func send(conn net.Conn, msg string, isUDP bool) {
	fmt.Printf("Sending \"%s\"\n> ", msg)

	encoder := encodeMsgTCP
	if isUDP {
		encoder = encodeMsgUDP
	}

	_, err := conn.Write(encoder(msg))
	if err != nil {
		log.Println(err)
	}
}

func readUserInput(out chan<- string) {
	defer close(out)

	userInput := bufio.NewScanner(os.Stdin)
	userInput.Split(bufio.ScanLines)

	fmt.Print("> ")
	for userInput.Scan() {
		text := userInput.Text()
		if text == "exit" {
			return
		}
		out <- text
	}
}

func listenForMessagesTCP(conn net.Conn, done chan struct{}) {
	defer close(done)

	scan.ListenAndDispatch(conn, onMessageReceivedTCP)
}

func onMessageReceivedTCP(prefix [scan.PrefixLengthBytes]byte, data []byte) {
	fmt.Println(string(data))
}

func encodeMsgTCP(s string) []byte {
	b := []byte(s)
	size := uint32(len(b))

	var prefix [scan.PrefixLengthBytes]byte
	binary.LittleEndian.PutUint32(prefix[:], size)

	out := make([]byte, 0, size+scan.PrefixLengthBytes)
	// Prefix
	out = append(out, prefix[:]...)
	// Data
	out = append(out, b...)
	return out
}

func listenForMessagesUDP(conn net.Conn, done chan struct{}) {
	defer close(done)

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		onMessageReceivedUDP(buf[:n])
	}
}

func onMessageReceivedUDP(data []byte) {
	fmt.Println(string(data))
}

func encodeMsgUDP(s string) []byte {
	return []byte(s)
}
