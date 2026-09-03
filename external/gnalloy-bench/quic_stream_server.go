package main

import (
	"context"
	"fmt"
	"sync"

	"gnalloy.org/transport-quic"
	"gnalloy.org/transport-quic/application"
)

const quicStreamALPNValue = "gnalloy-quic"

type quicStreamServer struct {
	addr     string
	listener quic.Listener
	payload  int
	codec    application.LengthPrefixedCodec
	buffers  sync.Pool
	executor *quicStreamExecutor
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func startQUICStreamServer(parent context.Context, cfg config) (*quicStreamServer, error) {
	if parent == nil {
		parent = context.Background()
	}
	tlsConfig, err := serverTLSConfig(cfg)
	if err != nil {
		return nil, err
	}
	listener, err := quic.ListenAddr(cfg.Addr, quic.Config{
		TLS:                   tlsConfig,
		NextProtos:            []string{quicStreamALPN(cfg)},
		MaxIncomingStreams:    maxIncomingQUICStreams(cfg.Connections),
		MaxIdleTimeout:        cfg.Timeout,
		KeepAlivePeriod:       cfg.Timeout / 4,
		InitialPacketSize:     quic.MinInitialPacketSize,
		MaxIncomingUniStreams: -1,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(parent)
	server := &quicStreamServer{
		addr:     listener.Addr().String(),
		listener: listener,
		payload:  cfg.Payload,
		codec:    application.LengthPrefixedCodec{MaxFrameSize: cfg.Payload},
		ctx:      ctx,
		cancel:   cancel,
	}
	server.buffers.New = func() any {
		buffer := make([]byte, cfg.Payload+2)
		return &buffer
	}
	server.executor = newQUICStreamExecutor(ctx, cfg.Connections, quicStreamQueueSize(cfg.Connections), server.serveStream)
	server.wg.Add(1)
	go server.accept()
	return server, nil
}

func (s *quicStreamServer) stop() {
	if s == nil {
		return
	}
	if s.cancel != nil {
		s.cancel()
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.wg.Wait()
	s.executor.stop()
}

func (s *quicStreamServer) accept() {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept(s.ctx)
		if err != nil {
			return
		}
		s.wg.Add(1)
		go s.serveConn(conn)
	}
}

func (s *quicStreamServer) serveConn(conn quic.Connection) {
	defer s.wg.Done()
	defer conn.CloseWithError(0, "benchmark done")
	for {
		stream, err := conn.AcceptStream(s.ctx)
		if err != nil {
			return
		}
		if !s.executor.submit(stream) {
			stream.CancelRead(0)
			stream.CancelWrite(0)
			return
		}
	}
}

func (s *quicStreamServer) serveStream(stream quic.Stream) {
	bufferPointer := s.buffers.Get().(*[]byte)
	buffer := (*bufferPointer)[:s.payload+2]
	defer s.buffers.Put(bufferPointer)
	payload, err := s.codec.ReadFrameInto(stream, buffer[2:])
	if err != nil {
		stream.CancelRead(0)
		stream.CancelWrite(0)
		return
	}
	if err := s.codec.WriteFrameInto(stream, payload, buffer); err != nil {
		stream.CancelWrite(0)
		return
	}
	_ = stream.Close()
}

func ensureQUICStreamConfig(cfg config) error {
	if cfg.TLSVersion != tlsVersion13 {
		return fmt.Errorf("%w: QUIC stream requires TLS 1.3", errInvalidConfig)
	}
	if quicStreamALPN(cfg) != quicStreamALPNValue {
		return fmt.Errorf("%w: QUIC stream requires ALPN %s", errInvalidConfig, quicStreamALPNValue)
	}
	if cfg.Payload > application.DefaultMaxFrameSize {
		return fmt.Errorf("%w: QUIC stream payload exceeds %d bytes", errInvalidConfig, application.DefaultMaxFrameSize)
	}
	return nil
}

func quicStreamALPN(cfg config) string {
	protocols := alpnProtocols(cfg.ALPN)
	if len(protocols) == 0 {
		return quicStreamALPNValue
	}
	return protocols[0]
}

func maxIncomingQUICStreams(connections int) int64 {
	streams := int64(connections) * 4
	if streams < 1024 {
		return 1024
	}
	return streams
}
