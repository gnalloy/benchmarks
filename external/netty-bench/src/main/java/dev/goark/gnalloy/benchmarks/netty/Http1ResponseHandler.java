package dev.goark.gnalloy.benchmarks.netty;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpResponseStatus;

import java.io.EOFException;
import java.io.IOException;
import java.time.Duration;
import java.util.Arrays;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.TimeUnit;

final class Http1ResponseHandler extends SimpleChannelInboundHandler<FullHttpResponse> {
    private final byte[] expected;
    private final byte[] reply;
    private final ArrayBlockingQueue<IOException> failures = new ArrayBlockingQueue<>(1);
    private final java.util.concurrent.Semaphore responses = new java.util.concurrent.Semaphore(0);

    Http1ResponseHandler(int payload) {
        expected = HttpPayload.body(payload);
        reply = new byte[payload];
    }

    @Override
    protected void channelRead0(ChannelHandlerContext context, FullHttpResponse response) {
        IOException failure = validate(response);
        if (failure != null) {
            failures.offer(failure);
        }
        responses.release();
    }

    void await(Duration timeout) throws IOException {
        try {
            if (!responses.tryAcquire(timeout.toNanos(), TimeUnit.NANOSECONDS)) {
                throw new IOException("netty-bench: HTTP/1 response timeout");
            }
        } catch (InterruptedException interrupted) {
            Thread.currentThread().interrupt();
            throw new IOException("netty-bench: interrupted while waiting for HTTP/1 response", interrupted);
        }
        IOException failure = failures.poll();
        if (failure != null) {
            throw failure;
        }
    }

    @Override
    public void channelInactive(ChannelHandlerContext context) throws Exception {
        failures.offer(new EOFException("netty-bench: connection closed"));
        responses.release();
        super.channelInactive(context);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext context, Throwable cause) {
        failures.offer(new IOException("netty-bench: HTTP/1 pipeline failure", cause));
        responses.release();
        context.close();
    }

    private IOException validate(FullHttpResponse response) {
        if (!HttpResponseStatus.OK.equals(response.status())) {
            return new IOException("netty-bench: unexpected status " + response.status());
        }
        ByteBuf content = response.content();
        if (content.readableBytes() != expected.length) {
            return new IOException("netty-bench: response length " + content.readableBytes()
                    + ", want " + expected.length);
        }
        content.getBytes(content.readerIndex(), reply);
        if (!Arrays.equals(reply, expected)) {
            return new IOException("netty-bench: response body mismatch");
        }
        return null;
    }
}
