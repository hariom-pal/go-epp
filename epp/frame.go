package epp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

const eppFrameHeaderLength = 4

const DefaultMaxFrameSize = 16 * 1024 * 1024

// ReadFrame reads a single RFC5734 EPP frame payload from conn.
func ReadFrame(conn net.Conn) ([]byte, error) {
	return ReadFrameWithMax(conn, DefaultMaxFrameSize)
}

// ReadFrameWithMax reads a single RFC5734 EPP frame payload using maxFrameSize.
func ReadFrameWithMax(conn net.Conn, maxFrameSize int) ([]byte, error) {
	header := make([]byte, eppFrameHeaderLength)

	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, classifyTransportError(err)
	}

	length := binary.BigEndian.Uint32(header)
	if length < eppFrameHeaderLength {
		return nil, newSDKError(
			ErrorKindFraming,
			fmt.Sprintf("invalid EPP frame length %d", length),
			ErrInvalidFrameLength,
		)
	}

	if maxFrameSize <= 0 {
		maxFrameSize = DefaultMaxFrameSize
	}
	if uint64(length) > uint64(maxFrameSize) {
		return nil, newSDKError(
			ErrorKindFraming,
			fmt.Sprintf("EPP frame length %d exceeds maximum %d", length, maxFrameSize),
			ErrFrameTooLarge,
		)
	}

	payload := make([]byte, length-eppFrameHeaderLength)

	if _, err := io.ReadFull(conn, payload); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			return nil, newSDKError(ErrorKindFraming, "truncated EPP frame", err)
		}
		return nil, classifyTransportError(err)
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
