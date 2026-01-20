package bufiox_test

import (
	"bufio"
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/nimelo/tui-go/bufiox"
)

type mockReader struct {
	data   []byte
	cursor int
}

func (m *mockReader) Read(ctx context.Context, p []byte) (int, error) {
	if m.cursor >= len(m.data) {
		return 0, io.EOF
	}

	n := copy(p, m.data[m.cursor:m.cursor+1])
	m.cursor += n

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return n, nil
	}
}
func TestContextScanner(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected [][]byte
	}{
		{
			name:  "Test1",
			input: []byte("one two three"),
			expected: [][]byte{
				[]byte("one"),
				[]byte("two"),
				[]byte("three"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			mockReader := &mockReader{data: test.input}

			scanner := bufiox.NewContextScanner(mockReader)
			scanner.Split(bufio.ScanWords)

			var results [][]byte
			for scanner.Scan(ctx) {
				token := scanner.Bytes()

				tmp := make([]byte, len(token))
				copy(tmp, token)

				results = append(results, tmp)
			}

			if !reflect.DeepEqual(results, test.expected) {
				t.Errorf("expected %q, got %q", test.expected, results)
			}
		})
	}
}
