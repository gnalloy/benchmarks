package dev.goark.gnalloy.benchmarks.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.LengthFieldPrepender;
import io.netty.handler.codec.quic.QuicChannel;
import io.netty.handler.codec.quic.QuicStreamChannel;
import io.netty.handler.codec.quic.QuicStreamType;
import io.netty.util.concurrent.Future;

import java.io.IOException;
import java.time.Duration;
import java.util.Arrays;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

record QuicStreamClient(Channel datagram, QuicChannel channel, byte[] payload, byte[] reply) implements AutoCloseable {
    void request(Duration timeout) throws Exception {
        ResponseHandler handler = new ResponseHandler(payload, reply);
        Future<QuicStreamChannel> streamFuture = channel.createStream(QuicStreamType.BIDIRECTIONAL, new ChannelInitializer<QuicStreamChannel>() {
            @Override
            protected void initChannel(QuicStreamChannel stream) {
                stream.pipeline()
                        .addLast(QuicStreamFrame.decoder())
                        .addLast(new LengthFieldPrepender(Short.BYTES))
                        .addLast(handler);
            }
        });
        QuicStreamChannel stream = streamFuture.get(timeout.toMillis(), TimeUnit.MILLISECONDS);
        ByteBuf request = stream.alloc().buffer(payload.length, payload.length).writeBytes(payload);
        stream.writeAndFlush(request)
                .addListener(QuicStreamChannel.SHUTDOWN_OUTPUT)
                .get(timeout.toMillis(), TimeUnit.MILLISECONDS);
        handler.await(timeout);
        stream.close().get(timeout.toMillis(), TimeUnit.MILLISECONDS);
    }

    String negotiatedProtocol() {
        String protocol = channel.sslEngine().getApplicationProtocol();
        return protocol == null ? "" : protocol;
    }

    @Override
    public void close() throws Exception {
        try {
            channel.close(true, 0, Unpooled.EMPTY_BUFFER).get(5, TimeUnit.SECONDS);
        } finally {
            datagram.close().sync();
        }
    }

    private static final class ResponseHandler extends SimpleChannelInboundHandler<ByteBuf> {
        private final byte[] expected;
        private final byte[] reply;
        private final CompletableFuture<Void> done = new CompletableFuture<>();

        private ResponseHandler(byte[] expected, byte[] reply) {
            this.expected = expected;
            this.reply = reply;
        }

        void await(Duration timeout) throws Exception {
            done.get(timeout.toMillis(), TimeUnit.MILLISECONDS);
        }

        @Override
        protected void channelRead0(ChannelHandlerContext ctx, ByteBuf payload) {
            try {
                if (payload.readableBytes() != reply.length) {
                    throw new IOException("netty-bench: quic stream response length mismatch");
                }
                payload.readBytes(reply);
                if (!Arrays.equals(reply, expected)) {
                    throw new IOException("netty-bench: quic stream response body mismatch");
                }
                done.complete(null);
                ctx.close();
            } catch (Throwable t) {
                done.completeExceptionally(t);
                ctx.close();
            }
        }

        @Override
        public void channelInactive(ChannelHandlerContext ctx) {
            if (!done.isDone()) {
                done.completeExceptionally(new IOException("netty-bench: quic stream closed"));
            }
        }

        @Override
        public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
            done.completeExceptionally(cause);
            ctx.close();
        }
    }
}
