package encoding

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	msg := []byte("This is a message that should be able to be encoded and decoded. " +
		"It should be able to be split and still maintain the ability to be decoded.")

	encoded, err := Encode([]byte(msg))
	if err != nil {
		t.Fatalf("failed to encode message: %v", err)
	}

	// wrap cipher text with padding
	padding := []byte("PaDdInGPaDdInGPaDdInGPaDdInGPaDdInGPaDdInG")
	data := buildSample(padding, encoded)

	buf := new(bytes.Buffer)
	dec := NewDecoder(buf)
	chunkSize := 100

	var i int
	for i = chunkSize; i < len(data); i += chunkSize {
		buf.Write(data[i-chunkSize : i])

		fmt.Printf("i: %d, len(msg):%d, len(padding):%d, len(data):%d",
			i, len(msg), len(padding), len(data))

		if i >= len(padding) && i <= len(data)-len(padding) {
			if !dec.More() {
				t.Fatalf("there should be a message being decoded at %d (%s)", i, string(data[:i]))
			} else if _, err = dec.Next(); err != ErrUnexpectedEOF {
				t.Fatalf("Did not return error '%v' as expected, instead returned: %v", ErrUnexpectedEOF, err)
			}
		} else if dec.More() {
			t.Fatalf("shouldn't be processing message")
		}
	}

	// write remaining data (if any left)
	if i < len(data) {
		buf.Write(data[i:])
	}

	actual, err := dec.Next()
	if err != nil {
		t.Fatalf("failed to decode end result: %v", err)
	}

	expected := buildSample(padding, msg)
	if string(actual) != string(expected) {
		t.Fatalf("expected (%s) did not match actual (%s)", string(expected), string(actual))
	}
}

func buildSample(padding, in []byte) []byte {
	data := append(padding, in...)
	data = append(data, padding...)
	return data
}
