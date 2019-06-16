package pgproxy

import (
	"log"

	"github.com/indent-go/apis/pkg/database/postgres/pgproto"
)

// processCopy determines how much of a COPY is remaining and when it's complete after the current messages.
func processCopy(msg []byte, copyLeft, length int, logger *log.Logger) (copyActive bool, left int) {
	copyActive = true
	left = copyLeft

	var msgType pgproto.MsgType
	var msgLen int32
	for i := 0; i < len(msg); i += int(msgLen + 1) {
		msgType = pgproto.Type(msg[i:])
		switch msgType {
		case pgproto.CopyDataMsg:
			msgLen = pgproto.Len(msg[i:])
			left = int(msgLen)
			logger.Printf("Received CopyData message (size: %d), conn to send data...", left)
		case pgproto.CopyDoneMsg:
			fallthrough
		case pgproto.CopyFailMsg:
			copyActive = false
			fallthrough
		default:
			// frame entirely of data
			left = left - length
		}

		if left >= length {
			left = left - length
			logger.Print("Accepting data...")
			return
		}
	}
	return
}

// copyMsg returns the type of COPY message.
func copyMsg(msg []byte) byte {
	return msg[0]
}
