package epp

import (
	"encoding/binary"
	"io"
	"net"
)

// ReadFrame reads a single RFC5734 EPP frame payload from conn.
func ReadFrame(conn net.Conn) ([]byte, error) {
	header := make([]byte, 4)

	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header)

	payload := make([]byte, length-4)

	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// WriteFrame writes payload as a single RFC5734 EPP frame to conn.
func WriteFrame(conn net.Conn, payload []byte) error {
	length := uint32(len(payload) + 4)

	frame := make([]byte, length)

	binary.BigEndian.PutUint32(frame[:4], length)

	copy(frame[4:], payload)

	_, err := conn.Write(frame)

	return err
}
