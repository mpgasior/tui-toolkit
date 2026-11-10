package seq

import "fmt"

type MalformedSeqErr struct {
	seqTyp string
	data   []byte
}

func (e *MalformedSeqErr) Error() string {
	return fmt.Sprintf("seq: malformed ansi %s sequence", e.seqTyp)
}

func (e *MalformedSeqErr) Data() []byte {
	return e.data
}

func errMalformedSeq(data []byte, seqType string) (int, []byte, error) {
	return 0, nil, &MalformedSeqErr{data: data, seqTyp: seqType}
}
