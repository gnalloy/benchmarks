package dev.goark.gnalloy.benchmarks.netty;

import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.handler.codec.http.HttpMethod;
import io.netty.handler.codec.http2.DefaultHttp2Headers;
import io.netty.handler.codec.http2.DefaultHttp2HeadersFrame;
import io.netty.handler.codec.http2.Http2Headers;
import io.netty.handler.codec.http2.Http2StreamChannel;
import io.netty.handler.codec.http2.Http2StreamChannelBootstrap;
import io.netty.util.concurrent.Future;

import java.io.IOException;
import java.time.Duration;
import java.util.concurrent.TimeUnit;

final class Http2Client implements AutoCloseable {
    private final Channel channel;
    private final Http2StreamChannelBootstrap streams;
    private final Http2Headers requestHeaders;
    private final byte[] expected;
    private final String negotiatedProtocol;

    Http2Client(Channel channel, String host, boolean tlsEnabled, int payload, String negotiatedProtocol) {
        this.channel = channel;
        this.streams = new Http2StreamChannelBootstrap(channel);
        this.requestHeaders = new DefaultHttp2Headers()
                .method(HttpMethod.GET.asciiName())
                .scheme(tlsEnabled ? "https" : "http")
                .path("/bench")
                .authority(host == null || host.isBlank() ? "127.0.0.1" : host);
        this.expected = HttpPayload.body(payload);
        this.negotiatedProtocol = negotiatedProtocol == null ? "" : negotiatedProtocol;
    }

    void exchange(Duration timeout) throws IOException {
        Http2ResponseHandler response = new Http2ResponseHandler(expected);
        Future<Http2StreamChannel> openFuture = streams.handler(response).open();
        if (!openFuture.awaitUninterruptibly(timeout.toNanos(), TimeUnit.NANOSECONDS)) {
            throw new IOException("netty-bench: HTTP/2 stream open timeout");
        }
        if (!openFuture.isSuccess()) {
            throw new IOException("netty-bench: HTTP/2 stream open failed", openFuture.cause());
        }
        Http2StreamChannel stream = openFuture.getNow();
        try {
            ChannelFuture writeFuture = stream.writeAndFlush(new DefaultHttp2HeadersFrame(requestHeaders, true));
            if (!writeFuture.awaitUninterruptibly(timeout.toNanos(), TimeUnit.NANOSECONDS)) {
                throw new IOException("netty-bench: HTTP/2 write timeout");
            }
            if (!writeFuture.isSuccess()) {
                throw new IOException("netty-bench: HTTP/2 write failed", writeFuture.cause());
            }
            response.await(timeout);
        } finally {
            stream.close().syncUninterruptibly();
        }
    }

    String negotiatedProtocol() {
        return negotiatedProtocol;
    }

    @Override
    public void close() {
        channel.close().syncUninterruptibly();
    }
}
