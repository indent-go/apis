package encoding

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestDecodeHeader(t *testing.T) {
	testMsg := []byte("This is a test message.")
	header := encodeHeader(testMsg)

	// base64 encode the header

	if length, err := decodeHeader(header); err != nil {
		t.Fatalf("failed to decode header: %v", err)
	} else if length != len(testMsg) {
		t.Fatalf("found incorrect length: expected '%d' and got '%d'", len(testMsg), length)
	}
}
