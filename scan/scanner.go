package scan

import (
	"bufio"
	"encoding/binary"
	"io"
)

// OnMessage receives the prefix bytes (data length), and the data itself when a message is received
type OnMessage func(prefix [PrefixLengthBytes]byte, data []byte)

// Listen listens on the given Reader using a Scanner, for length-prefixed messages.
func ListenAndDispatch(r io.Reader, dispatch OnMessage) {
	scanner := bufio.NewScanner(r)
	lp := LengthPrefixer{}
	scanner.Split(lp.ScanLengthPrefix)

	var prefix [PrefixLengthBytes]byte
	for scanner.Scan() {
		b := scanner.Bytes()

		// Send the message, including the length prefix
		// This needs to happen in one "send", or we'll end up seeing stuff like [length, length, data, data]
		switch lp.State {
		case ReadData:
			lp.State = ReadPrefix
			binary.LittleEndian.PutUint32(prefix[:], lp.Length)
			dispatch(prefix, b)
		default: // ReadPrefix
			lp.State = ReadData
			lp.Length = binary.LittleEndian.Uint32(b)
		}
	}
}
