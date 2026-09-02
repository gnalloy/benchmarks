package httpbench

import (
	"context"

	"gnalloy.org/codec-http1"
	"gnalloy.org/gnalloy/channel"
	gnalloytls "gnalloy.org/handler-tls"
)

const clientReadBufferSize = 16 * 1024

type protocolClient struct {
	channel  channel.Channel
	response *responseHandler
	request  http1.Request
}

func newProtocolClient(host string, expectedBody []byte) *protocolClient {
	return &protocolClient{
		response: newResponseHandler(expectedBody),
		request: http1.Request{
			Method: "GET",
			URI:    "/bench",
			Headers: http1.Headers{
				"Host":       host,
				"User-Agent": "gnalloy-bench",
				"Accept":     "*/*",
				"Connection": "keep-alive",
			},
		},
	}
}

func (c *protocolClient) addPipeline(ch channel.Channel, cfg Config) error {
	if cfg.TLS != nil {
		tlsConfig := cfg.TLS.Clone()
		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName = cfg.ServerName
		}
		if err := ch.Pipeline().AddLast("tls", gnalloytls.Client(gnalloytls.Config{TLS: tlsConfig})); err != nil {
			return err
		}
	}
	decoder, err := http1.NewResponseDecoder(16*1024, len(c.response.expectedBody))
	if err != nil {
		return err
	}
	if err := ch.Pipeline().AddLast("httpEncoder", http1.NewRequestEncoder()); err != nil {
		return err
	}
	if err := ch.Pipeline().AddLast("httpDecoder", decoder); err != nil {
		return err
	}
	return ch.Pipeline().AddLast("response", c.response)
}

func (c *protocolClient) exchange(ctx context.Context) error {
	if err := c.channel.WriteAndFlush(c.request); err != nil {
		return err
	}
	select {
	case err := <-c.response.responses:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
