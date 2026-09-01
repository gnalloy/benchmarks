package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	quicStreamALPNValue    = "gnalloy-quic"
	quicStreamMaxFrameSize = 0xffff
)

// writeQUICStreamFrame 使用两字节网络序长度前缀写出应用消息。
func writeQUICStreamFrame(w io.Writer, payload []byte) error {
	if w == nil {
		return errInvalidConfig
	}
	if len(payload) > quicStreamMaxFrameSize || len(payload) > 0xffff {
		return fmt.Errorf("%w: QUIC stream frame %d exceeds %d bytes", errInvalidConfig, len(payload), quicStreamMaxFrameSize)
	}
	var header [2]byte
	binary.BigEndian.PutUint16(header[:], uint16(len(payload)))
	if err := writeAll(w, header[:]); err != nil {
		return err
	}
	return writeAll(w, payload)
}

// readQUICStreamFrame 将一帧应用消息读入调用方复用的缓冲区。
func readQUICStreamFrame(r io.Reader, dst []byte, maxFrameSize int) ([]byte, error) {
	if r == nil {
		return nil, errInvalidConfig
	}
	if maxFrameSize <= 0 {
		maxFrameSize = quicStreamMaxFrameSize
	}
	var header [2]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}
	size := int(binary.BigEndian.Uint16(header[:]))
	if size > maxFrameSize || size > len(dst) {
		return nil, fmt.Errorf("%w: QUIC stream frame %d exceeds buffer %d or max %d", errInvalidConfig, size, len(dst), maxFrameSize)
	}
	if _, err := io.ReadFull(r, dst[:size]); err != nil {
		return nil, err
	}
	return dst[:size], nil
}
