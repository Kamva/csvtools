package csvtools

import (
	"bufio"
	"bytes"
)

// ScanGroupedCSVRecords is a scanner that you can use as bufio.Scanner's splitter.
// e.g.,
// sc:=bufio.NewScanner()
// var split ScanGroupedCSVRecords
// sc.Split(split.SplitFunc)
type ScanGroupedCSVRecords struct {
	head []byte
	prev []byte
}

func (s *ScanGroupedCSVRecords) SplitFunc(data []byte, efo bool) (advance int, token []byte, err error) {
	if efo && len(data) == 0 {
		return 0, nil, nil
	}

	if s.head == nil { // catch head for all tokens and drop it.
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			s.head=make([]byte,i+1)
			copy(s.head,dropCR(data[0 : i+1])) // for head we keep new line too.
			// We have a full newline-terminated line.
			return i + 1, nil, nil
		}
	}

	if s.prev == nil { // initialize for first record.
		if i := bytes.IndexByte(data, ','); i >= 0 {
			s.prev = dropCR(data[0:i])
		}
	}
	var i = 0
	for {
		newI := bytes.IndexByte(data[i:], '\n')
		if newI<0{
			break // break to go out of for and check efo or get more.
		}
		// We have a full newline-terminated line.
		// so we should check our csc record's key.
		j := bytes.IndexByte(data[i+newI:], ',')

		if j < 0 {
			break // break to go out of for and check efo or get more.
		}
		i = i + newI

		old := s.prev
		s.prev = dropCR(data[i+1 : i+j])
		if !bytes.Equal(old, s.prev) { // If we don't have the same keys, return token.
			return i + 1, append(s.head, dropCR(data[0:i])...), nil // return head + record.
		}

		i = i + 1
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if efo {
		return len(data), append(s.head,dropCR(data)...), nil
	}
	// Request more data.
	return 0, nil, nil
}

var _ bufio.SplitFunc = (&ScanGroupedCSVRecords{}).SplitFunc
