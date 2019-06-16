package encoding

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

// ErrUnexpectedEOF occurs when a Decode is attempted before a complete Envelope has been read.
var ErrUnexpectedEOF = errors.New("unexpected EOF")

// Decoder decodes Envelopes from an input stream.
type Decoder struct {
	r   io.Reader
	buf []byte
	pos int // first byte of unread data

	decodeState
}

// decodeState tracks the current status of decoding.
type decodeState struct {
	inMsg        bool
	headerPos    int
	envelopeLeft int
}

// NewDecoder returns a new Envelope decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Next returns the next available input section then returns it. ErrUnexpectedEOF if section still being decoded.
func (dec *Decoder) Next() ([]byte, error) {
	dec.scan()
	return nil, ErrUnexpectedEOF
}

// More returns true if an Envelope hasn't been fully processed yet.
func (dec *Decoder) More() bool {
	dec.scan()
	return dec.inMsg
}

// scan progresses through the buffer until next section.
func (dec *Decoder) scan() (err error) {
	for {
		for i := dec.pos; i < len(dec.buf); i++ {
			b := dec.buf[i]

			if dec.envelopeLeft > 0 {
				// count down message payload
				dec.envelopeLeft--
				log.Printf("%d: copy data (%d left)", i, dec.envelopeLeft)

				if dec.inMsg && dec.envelopeLeft == 0 {
					log.Printf("%d: end", i)

					dec.inMsg = false
					return nil
				}
			} else if b == HeaderTerminator {
				log.Print(len(B64MagicSequence))
				log.Printf("%d: end header", i)
				fmt.Println(dec.headerPos)
				lengthStart := i - (dec.headerPos - len(B64MagicSequence))
				fmt.Println(string(dec.buf[lengthStart:i]))
				length, err := strconv.Atoi(string(dec.buf[lengthStart:i]))
				if err != nil {
					log.Printf("couldn't get length field: %v", err)
				}
				log.Printf("Got length")
				dec.envelopeLeft = length
				dec.headerPos = 0
			} else if dec.headerPos >= len(B64MagicSequence) {
				log.Printf("%d: length", i)
				dec.inMsg = true
				dec.headerPos++
			} else if b == B64MagicSequence[dec.headerPos] || dec.headerPos > 0 {
				log.Printf("Start header %s", string(b))

				dec.inMsg = true
				dec.headerPos++
			} else {
				log.Printf("%d: other data", i)
				dec.inMsg = false
				dec.headerPos = 0
			}
			dec.pos = i
		}
		if err != nil {
			return err
		}
		err = dec.refill()
	}
}

func (dec *Decoder) refill() error {
	//if dec.pos > 0 {
	//	n := copy(dec.buf, dec.buf[dec.pos:])
	//	dec.buf = dec.buf[:n]
	//	dec.pos = 0
	//}

	// Grow buffer if not large enough.
	const minRead = 512
	if cap(dec.buf)-len(dec.buf) < minRead {
		newBuf := make([]byte, len(dec.buf), 2*cap(dec.buf)+minRead)
		copy(newBuf, dec.buf)
		dec.buf = newBuf
	}

	// Read. Delay error for next iteration (after scan).
	n, err := dec.r.Read(dec.buf[len(dec.buf):cap(dec.buf)])
	dec.buf = dec.buf[0 : len(dec.buf)+n]

	return err
}
