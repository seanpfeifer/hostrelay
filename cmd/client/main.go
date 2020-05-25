package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/seanpfeifer/hostrelay/scan"
)

const (
	host = ":8080"
)

func main() {
	done := make(chan struct{})
	out := make(chan string)

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go listenForMessages(conn, done)
	go readUserInput(out)

	for {
		select {
		case <-done:
			return
		case msg, ok := <-out:
			if !ok {
				return
			}
			send(conn, msg)
		}
	}

}

func send(conn net.Conn, msg string) {
	fmt.Printf("Sending \"%s\"\n> ", msg)

	_, err := conn.Write(encodeMsg(msg))
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

func listenForMessages(conn net.Conn, done chan struct{}) {
	defer close(done)

	scan.ListenAndDispatch(conn, onMessageReceived)
}

func onMessageReceived(prefix [scan.PrefixLengthBytes]byte, data []byte) {
	fmt.Println(string(data))
}

func encodeMsg(s string) []byte {
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
