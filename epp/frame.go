package epp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const eppFrameHeaderLength = 4

// ReadFrame reads a single RFC5734 EPP frame payload from conn.
func ReadFrame(conn net.Conn) ([]byte, error) {
	header := make([]byte, eppFrameHeaderLength)

	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header)
	if length < eppFrameHeaderLength {
		return nil, fmt.Errorf("invalid EPP frame length %d", length)
	}

	payload := make([]byte, length-eppFrameHeaderLength)

	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// WriteFrame writes payload as a single RFC5734 EPP frame to conn.
func WriteFrame(conn net.Conn, payload []byte) error {
	maxPayloadLength := uint64(^uint32(0)) - eppFrameHeaderLength
	if uint64(len(payload)) > maxPayloadLength {
		return fmt.Errorf("EPP frame payload too large: %d bytes", len(payload))
	}

	length := uint32(len(payload) + eppFrameHeaderLength)

	frame := make([]byte, length)

	binary.BigEndian.PutUint32(frame[:eppFrameHeaderLength], length)

	copy(frame[eppFrameHeaderLength:], payload)

	for len(frame) > 0 {
		n, err := conn.Write(frame)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		frame = frame[n:]
	}

	return nil
}
