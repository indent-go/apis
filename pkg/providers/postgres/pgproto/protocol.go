package pgproto

import (
	"bytes"
	"encoding/binary"
)

// Protocol field positions
const (
	// MsgTypePos is the position in a message indicating it's type.
	MsgTypePos = 0
	// MsgLenStart is the position that starts the length field.
	MsgLenStart = 1
	// MsgLenEnd is the position that ends the length field.
	MsgLenEnd = 5
)

// Type of the given msg.
func Type(msg []byte) MsgType {
	return MsgType(msg[MsgTypePos])
}

// QueryBody for a msg type of Query.
func QueryBody(msg []byte) []byte {
	len := Len(msg)
	return msg[MsgLenEnd:len]
}

// Len returns length of a given msg including the length field itself.
func Len(msg []byte) (length int32) {
	r := bytes.NewReader(msg[MsgLenStart:MsgLenEnd])
	binary.Read(r, binary.BigEndian, &length)
	return length
}
