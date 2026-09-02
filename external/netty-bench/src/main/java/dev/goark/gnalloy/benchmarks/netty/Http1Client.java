package dev.goark.gnalloy.benchmarks.netty;

import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.handler.codec.http.DefaultHttpRequest;
import io.netty.handler.codec.http.HttpHeaderNames;
import io.netty.handler.codec.http.HttpHeaderValues;
import io.netty.handler.codec.http.HttpMethod;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpVersion;
import io.netty.handler.codec.http.LastHttpContent;

import java.io.IOException;
import java.time.Duration;

final class Http1Client implements AutoCloseable {
    private final Channel channel;
    private final Http1ResponseHandler responses;
    private final HttpRequest request;
    private final String negotiatedProtocol;

    Http1Client(Channel channel, Http1ResponseHandler responses, String host, String negotiatedProtocol) {
        this.channel = channel;
        this.responses = responses;
        this.request = new DefaultHttpRequest(HttpVersion.HTTP_1_1, HttpMethod.GET, "/bench");
        request.headers().set(HttpHeaderNames.HOST, host);
        request.headers().set(HttpHeaderNames.USER_AGENT, "netty-bench");
        request.headers().set(HttpHeaderNames.ACCEPT, "*/*");
        request.headers().set(HttpHeaderNames.CONNECTION, HttpHeaderValues.KEEP_ALIVE);
        this.negotiatedProtocol = negotiatedProtocol == null ? "" : negotiatedProtocol;
    }

    void exchange(Duration timeout) throws IOException {
        channel.write(request);
        ChannelFuture future = channel.writeAndFlush(LastHttpContent.EMPTY_LAST_CONTENT);
        if (!future.awaitUninterruptibly(timeout.toNanos(), java.util.concurrent.TimeUnit.NANOSECONDS)) {
            throw new IOException("netty-bench: HTTP/1 write timeout");
        }
        if (!future.isSuccess()) {
            throw new IOException("netty-bench: HTTP/1 write failed", future.cause());
        }
        responses.await(timeout);
    }

    String negotiatedProtocol() {
        return negotiatedProtocol;
    }

    @Override
    public void close() {
        channel.close().syncUninterruptibly();
    }
}
