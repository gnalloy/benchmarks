package benchh2

import (
	http2 "gnalloy.org/codec-http2"
	"gnalloy.org/gnalloy/channel"
	gnalloytls "gnalloy.org/handler-tls"
)

func addClientPipeline(ch channel.Channel, cfg Config, response *responseHandler) error {
	pipeline := ch.Pipeline()
	if cfg.TLS != nil {
		if err := pipeline.AddLast("tls", gnalloytls.Client(gnalloytls.Config{TLS: cfg.TLS})); err != nil {
			return err
		}
	}
	frameDecoder, err := http2.NewFrameDecoder(http2.DefaultMaxFrameSize)
	if err != nil {
		return err
	}
	headerDecoder, err := http2.NewHeaderDecoder(http2.HeaderCodecConfig{})
	if err != nil {
		return err
	}
	headerEncoder, err := http2.NewHeaderEncoder(http2.HeaderCodecConfig{})
	if err != nil {
		return err
	}
	multiplexer, err := http2.NewStreamMultiplexer(http2.MultiplexerConfig{Server: false})
	if err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-preface", http2.NewPrefaceEncoder()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-frame-decoder", frameDecoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-typed-decoder", http2.NewTypedFrameDecoder()); err != nil {
		return err
	}
	frameEncoder := http2.NewFrameEncoder()
	if cfg.TLS != nil {
		frameEncoder = http2.NewFrameEncoderWithConfig(http2.FrameEncoderConfig{CoalescePayload: true})
	}
	if err := pipeline.AddLast("http2-frame-encoder", frameEncoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-typed-encoder", http2.NewTypedFrameEncoder()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-header-decoder", headerDecoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-header-encoder", headerEncoder); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-settings-ack", http2.NewSettingsAckHandler()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-ping-ack", http2.NewPingAckHandler()); err != nil {
		return err
	}
	if err := pipeline.AddLast("http2-multiplexer", multiplexer); err != nil {
		return err
	}
	return pipeline.AddLast("http2-response", response)
}
