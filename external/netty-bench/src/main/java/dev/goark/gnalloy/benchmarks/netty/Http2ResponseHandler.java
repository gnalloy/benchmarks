package dev.goark.gnalloy.benchmarks.netty;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http2.Http2DataFrame;
import io.netty.handler.codec.http2.Http2HeadersFrame;
import io.netty.handler.codec.http2.Http2StreamFrame;

import java.io.EOFException;
import java.io.IOException;
import java.time.Duration;
import java.util.Arrays;
import java.util.concurrent.Semaphore;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

final class Http2ResponseHandler extends SimpleChannelInboundHandler<Http2StreamFrame> {
    private final byte[] expected;
    private final byte[] reply;
    private final Semaphore completed = new Semaphore(0);
    private final AtomicBoolean completionSignaled = new AtomicBoolean();
    private volatile IOException failure;
    private int received;
    private boolean statusOK;

    Http2ResponseHandler(byte[] expected) {
        this.expected = expected;
        this.reply = new byte[expected.length];
    }

    @Override
    protected void channelRead0(ChannelHandlerContext context, Http2StreamFrame frame) {
        if (frame instanceof Http2HeadersFrame headersFrame) {
            if (headersFrame.headers().status() != null) {
                statusOK = HttpResponseStatus.OK.codeAsText().equals(headersFrame.headers().status());
            }
            if (headersFrame.isEndStream()) {
                complete(validate());
            }
            return;
        }
        if (frame instanceof Http2DataFrame dataFrame) {
            readData(dataFrame);
        }
    }

    void await(Duration timeout) throws IOException {
        try {
            if (!completed.tryAcquire(timeout.toNanos(), TimeUnit.NANOSECONDS)) {
                throw new IOException("netty-bench: HTTP/2 response timeout");
            }
        } catch (InterruptedException interrupted) {
            Thread.currentThread().interrupt();
            throw new IOException("netty-bench: interrupted while waiting for HTTP/2 response", interrupted);
        }
        if (failure != null) {
            throw failure;
        }
    }

    @Override
    public void channelInactive(ChannelHandlerContext context) throws Exception {
        complete(new EOFException("netty-bench: HTTP/2 stream closed before response completed"));
        super.channelInactive(context);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext context, Throwable cause) {
        complete(new IOException("netty-bench: HTTP/2 pipeline failure", cause));
        context.close();
    }

    private void readData(Http2DataFrame frame) {
        ByteBuf content = frame.content();
        int readable = content.readableBytes();
        if (readable > reply.length - received) {
            complete(new IOException("netty-bench: HTTP/2 response body too large"));
            return;
        }
        content.getBytes(content.readerIndex(), reply, received, readable);
        received += readable;
        if (frame.isEndStream()) {
            complete(validate());
        }
    }

    private IOException validate() {
        if (!statusOK) {
            return new IOException("netty-bench: HTTP/2 response status is not 200");
        }
        if (received != expected.length) {
            return new IOException("netty-bench: HTTP/2 response body length " + received + ", want " + expected.length);
        }
        if (!Arrays.equals(reply, expected)) {
            return new IOException("netty-bench: HTTP/2 response body mismatch");
        }
        return null;
    }

    private void complete(IOException error) {
        if (!completionSignaled.compareAndSet(false, true)) {
            return;
        }
        failure = error;
        completed.release();
    }
}
