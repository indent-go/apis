package pgproxy

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"testing"

	"go.indent.com/apis/pkg/providers/postgres/pgproto"
)

func TestCopyOneMsgRecieve(t *testing.T) {
	copyData := createMessage(pgproto.CopyDataMsg, []byte("test msg")...)
	copyDone := createMessage(pgproto.CopyDoneMsg)

	logger := log.New(os.Stderr, "", log.LstdFlags)
	data := append(copyData, copyDone...)

	copyActive, _ := processCopy(data, 0, len(data), logger)
	if copyActive {
		t.Fatalf("Should have stopped COPY mode.")
	}
}

func TestCopyTwoMsgReceive(t *testing.T) {
	copyData := createMessage(pgproto.CopyDataMsg, []byte("test msg")...)

	logger := log.New(os.Stderr, "", log.LstdFlags)
	copyActive, left := processCopy(copyData, 0, len(copyData), logger)
	if !copyActive {
		t.Fatalf("Should have started COPY mode.")
	}

	copyDone := createMessage(pgproto.CopyDoneMsg)
	copyActive, left = processCopy(copyDone, left, len(copyDone), logger)
	if copyActive {
		t.Fatalf("Should have stopped COPY mode.")
	}
}

func TestCopy(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	copyData := createMessage(pgproto.CopyDataMsg, []byte("This message tests when it is split up into pieces")...)
	copyActive, left := processCopy(copyData[:9], 0, len(copyData[:9]), logger)
	if !copyActive {
		t.Fatalf("Should have started COPY mode.")
	}

	end := append(copyData[10:], createMessage(pgproto.CopyDoneMsg)...)
	copyActive, left = processCopy(end, left, len(end), logger)
	if copyActive {
		t.Fatalf("Should have stopped COPY mode.")
	}
}

func createMessage(t pgproto.MsgType, payload ...byte) []byte {
	// compensate for length + payload + null
	length := 5 + len(payload) + 1

	// write message type
	buf := bytes.NewBuffer([]byte{byte(t)})
	buf.Grow(length)

	// write length
	lengthField := make([]byte, 5)
	binary.BigEndian.PutUint32(lengthField, uint32(length))
	buf.Write(lengthField)

	// write payload
	buf.Write(payload)

	// write null
	buf.Write([]byte{0x00})

	return buf.Bytes()
}
