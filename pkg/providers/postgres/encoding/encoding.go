package encoding

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// MagicSequence is used to determine the beginning of an encoded Envelope.
	MagicSequence = "\xF0\x9F\x94\x92"

	// HeaderTerminator is the character used mark the end of a encodeHeader.
	HeaderTerminator = 'D'
)

var (
	// B64Encoding is the base64 variant used to encode data.
	B64Encoding = base64.RawURLEncoding

	// ByteMagicSequence is the MagicSequence in []byte form.
	ByteMagicSequence = []byte(MagicSequence)

	// B64MagicSequence is the base64 encoding of the MagicSequence.
	B64MagicSequence = B64Encoding.EncodeToString(ByteMagicSequence)
)

var (
	ErrNoEnvelopeInText = errors.New("no envelope in text")
	ErrNoHeaderEnd      = errors.New("the end of the header could not be found")
)

// Encode produces a format which can be distributed through a variety of methods.
func Encode(payload []byte) ([]byte, error) {
	// base64 encode msg
	msg := make([]byte, B64Encoding.EncodedLen(len(payload)))
	B64Encoding.Encode(msg, payload)

	// create encodeHeader and append data to it
	out := append(encodeHeader(msg), msg...)
	return out, nil
}

// encodeHeader generated using MagicSequence and content length.
func encodeHeader(data []byte) []byte {
	length := strconv.Itoa(len(data))
	h := []byte(B64MagicSequence + length)
	return append(h, HeaderTerminator)
}

// Decode restores a representation produced with Encode.
func Decode(content []byte) ([]byte, error) {
	// search for end of header
	hEnd := 0
	for i := 0; i < len(content); i++ {
		if content[i] == HeaderTerminator {
			hEnd = i
			break
		}
	}

	if hEnd == 0 {
		return nil, errors.New("could not find header")
	}

	// strip header
	content = content[hEnd+1:]

	data := make([]byte, B64Encoding.DecodedLen(len(content)))
	if _, err := B64Encoding.Decode(data, content); err != nil {
		return nil, fmt.Errorf("could not decode content: %v", err)
	}
	return data, nil
}

// decodeHeader validates the encodeHeader and returns the size of the message.
func decodeHeader(header []byte) (int, error) {
	strHeader := string(header)

	// validate encodeHeader
	minLength := len(B64MagicSequence) + 2
	if len(strHeader) < minLength {
		return 0, fmt.Errorf("header is only %d, must be at least %d", len(strHeader), minLength)
	} else if strHeader[0:len(B64MagicSequence)] != B64MagicSequence {
		return 0, fmt.Errorf("header must start with '%s', header: '%s'", B64MagicSequence, strHeader)
	} else if strHeader[len(strHeader)-1] != HeaderTerminator {
		return 0, fmt.Errorf("header must end with '%x', header: '%s'", HeaderTerminator, strHeader)
	}

	// decode length
	lengthStr := strHeader[len(B64MagicSequence) : len(strHeader)-1]
	if msgLen, err := strconv.Atoi(lengthStr); err != nil {
		return 0, fmt.Errorf("could not decode length from specifier '%s': %v", lengthStr, err)
	} else {
		return msgLen, nil
	}
}

// DecodeText detects messages and returns the text with the expanded Envelope.
func DecodeText(text string) (result string, start, end int, err error) {
	// determine if sequence includes Envelope
	start = strings.Index(text, B64MagicSequence)
	if start == -1 {
		err = ErrNoEnvelopeInText
		return
	}

	// detect position of header
	headerEnd := -1
	for i := start + len(B64MagicSequence); i < len(text); i++ {
		if text[i] == HeaderTerminator {
			headerEnd = i + 1
			break
		}
	}

	if headerEnd == -1 {
		err = ErrNoHeaderEnd
		return
	}

	// decode header to determine envelope length
	envelopeLen, err := decodeHeader([]byte(text[start:headerEnd]))
	if err != nil {
		err = fmt.Errorf("failed to decode header: %v", err)
		return
	}
	// decode envelope
	end = headerEnd + envelopeLen

	decoded, err := Decode([]byte(text[start:end]))
	if err != nil {
		err = fmt.Errorf("failed to decode envelope: %v", err)
		return
	}
	result = string(decoded)
	return
}
