package dev.goark.gnalloy.benchmarks.netty;

import io.netty.bootstrap.Bootstrap;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.quic.InsecureQuicTokenHandler;
import io.netty.handler.codec.quic.QuicServerCodecBuilder;
import io.netty.handler.codec.quic.QuicSslContext;
import io.netty.handler.codec.quic.QuicSslContextBuilder;
import io.netty.handler.codec.quic.QuicStreamChannel;
import io.netty.handler.ssl.util.SelfSignedCertificate;

import java.net.InetSocketAddress;
import java.util.concurrent.TimeUnit;

final class QuicStreamServer implements AutoCloseable {
    private static final long STREAM_WINDOW = 6L * 1024 * 1024;
    private static final long CONNECTION_WINDOW = 16L * 1024 * 1024;

    private final DatagramEventLoopResources resources;
    private final Channel channel;

    private QuicStreamServer(DatagramEventLoopResources resources, Channel channel) {
        this.resources = resources;
        this.channel = channel;
    }

    static QuicStreamServer start(Config config) throws Exception {
        DatagramEventLoopResources resources = DatagramEventLoopResources.create(config);
        try {
            ChannelHandler codec = serverCodec(config);
            Channel channel = new Bootstrap()
                    .group(resources.group())
                    .channel(resources.channelType())
                    .option(ChannelOption.SO_REUSEADDR, true)
                    .option(ChannelOption.ALLOCATOR, PooledByteBufAllocator.DEFAULT)
                    .handler(codec)
                    .bind(config.host(), config.port())
                    .sync()
                    .channel();
            return new QuicStreamServer(resources, channel);
        } catch (Throwable t) {
            resources.close();
            if (t instanceof InterruptedException interrupted) {
                throw interrupted;
            }
            if (t instanceof Exception exception) {
                throw exception;
            }
            if (t instanceof RuntimeException runtimeException) {
                throw runtimeException;
            }
            throw new RuntimeException(t);
        }
    }

    InetSocketAddress address() {
        return (InetSocketAddress) channel.localAddress();
    }

    @Override
    public void close() throws InterruptedException {
        ChannelFuture closeFuture = channel.close().sync();
        closeFuture.await();
        resources.close();
    }

    private static ChannelHandler serverCodec(Config config) throws Exception {
        SelfSignedCertificate certificate = new SelfSignedCertificate("gnalloy.local");
        QuicSslContext sslContext = QuicSslContextBuilder
                .forServer(certificate.key(), null, certificate.cert())
                .applicationProtocols(config.alpnProtocols().toArray(String[]::new))
                .build();
        return new QuicServerCodecBuilder()
                .sslContext(sslContext)
                .maxIdleTimeout(config.timeout().toMillis(), TimeUnit.MILLISECONDS)
                .initialMaxData(CONNECTION_WINDOW)
                .initialMaxStreamDataBidirectionalLocal(STREAM_WINDOW)
                .initialMaxStreamDataBidirectionalRemote(STREAM_WINDOW)
                .initialMaxStreamDataUnidirectional(STREAM_WINDOW)
                .initialMaxStreamsBidirectional(QuicLimits.bidirectionalStreamLimit(config))
                .initialMaxStreamsUnidirectional(0)
                .tokenHandler(InsecureQuicTokenHandler.INSTANCE)
                .handler(new ChannelInitializer<Channel>() {
                    @Override
                    protected void initChannel(Channel channel) {
                        // 每条 QUIC 连接使用独立 pipeline，避免复用非 sharable handler。
                    }
                })
                .streamHandler(new ChannelInitializer<QuicStreamChannel>() {
                    @Override
                    protected void initChannel(QuicStreamChannel channel) {
                        channel.pipeline()
                                .addLast(QuicStreamFrame.decoder())
                                .addLast(new EchoHandler());
                    }
                })
                .build();
    }

    private static final class EchoHandler extends SimpleChannelInboundHandler<io.netty.buffer.ByteBuf> {
        @Override
        protected void channelRead0(ChannelHandlerContext ctx, io.netty.buffer.ByteBuf payload) {
            ctx.writeAndFlush(QuicStreamFrame.encode(ctx.alloc(), payload))
                    .addListener(QuicStreamChannel.SHUTDOWN_OUTPUT)
                    .addListener(ChannelFutureListener.CLOSE);
        }

        @Override
        public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
            ctx.close();
        }
    }
}
