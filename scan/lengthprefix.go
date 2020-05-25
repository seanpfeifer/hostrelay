package scan

const (
	PrefixLengthBytes            = 4
	ReadPrefix        ReadyState = false
	ReadData          ReadyState = true
)

type ReadyState bool

type LengthPrefixer struct {
	Length uint32
	State  ReadyState // If we're ready to read data out of the prefixer
}

// ScanLengthPrefix is implemented similar to bufio.ScanLines, except using the length of data read so far as the delimiter
func (lp *LengthPrefixer) ScanLengthPrefix(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	switch lp.State {
	case ReadData:
		// Slice out the length needed for the message
		if len(data) >= int(lp.Length) {
			return int(lp.Length), data[:lp.Length], nil
		}
	default: // ReadPrefix
		// Slice out the length needed for the prefix
		if len(data) >= PrefixLengthBytes {
			return PrefixLengthBytes, data[:PrefixLengthBytes], nil
		}
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
