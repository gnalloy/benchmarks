package dev.goark.gnalloy.benchmarks.netty;

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.epoll.Epoll;
import io.netty.channel.epoll.EpollEventLoopGroup;
import io.netty.channel.epoll.EpollSocketChannel;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.handler.codec.http2.Http2FrameCodecBuilder;
import io.netty.handler.codec.http2.Http2MultiplexHandler;
import io.netty.handler.codec.http2.Http2Settings;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslHandler;

import java.net.InetSocketAddress;
import java.util.concurrent.TimeUnit;

final class Http2ClientGroup implements AutoCloseable {
    private final EventLoopGroup eventLoops;
    private final Http2Client[] clients;

    private Http2ClientGroup(EventLoopGroup eventLoops, Http2Client[] clients) {
        this.eventLoops = eventLoops;
        this.clients = clients;
    }

    static Http2ClientGroup connect(InetSocketAddress address, Config config) throws Exception {
        EventLoopGroup eventLoops = newEventLoopGroup(config);
        Http2Client[] clients = new Http2Client[config.connections()];
        try {
            SslContext sslContext = SslSupport.clientContext(config);
            Bootstrap bootstrap = new Bootstrap()
                    .group(eventLoops)
                    .channel(clientChannelType(config))
                    .option(ChannelOption.TCP_NODELAY, true)
                    .option(ChannelOption.CONNECT_TIMEOUT_MILLIS, Math.toIntExact(config.timeout().toMillis()))
                    .handler(new ChannelInitializer<SocketChannel>() {
                        @Override
                        protected void initChannel(SocketChannel channel) {
                            if (sslContext != null) {
                                channel.pipeline().addLast(sslContext.newHandler(
                                        channel.alloc(), SslSupport.SERVER_NAME, address.getPort()));
                            }
                            Http2Settings settings = new Http2Settings()
                                    .pushEnabled(false)
                                    .initialWindowSize(1 << 30);
                            channel.pipeline().addLast(Http2FrameCodecBuilder.forClient()
                                    .initialSettings(settings)
                                    .autoAckSettingsFrame(true)
                                    .autoAckPingFrame(true)
                                    .build());
                            channel.pipeline().addLast(new Http2MultiplexHandler(Http2InboundStreamHandler.INSTANCE));
                        }
                    });
            for (int i = 0; i < clients.length; i++) {
                Channel channel = bootstrap.connect(address).sync().channel();
                SslHandler sslHandler = channel.pipeline().get(SslHandler.class);
                String negotiatedProtocol = "";
                if (sslHandler != null) {
                    sslHandler.handshakeFuture().sync();
                    negotiatedProtocol = sslHandler.applicationProtocol();
                }
                clients[i] = new Http2Client(
                        channel,
                        config.host(),
                        config.tlsEnabled(),
                        config.payload(),
                        negotiatedProtocol);
            }
            return new Http2ClientGroup(eventLoops, clients);
        } catch (Throwable failure) {
            closeClients(clients);
            eventLoops.shutdownGracefully(0, 0, TimeUnit.MILLISECONDS).syncUninterruptibly();
            if (failure instanceof Exception exception) {
                throw exception;
            }
            throw new RuntimeException(failure);
        }
    }

    Http2Client[] clients() {
        return clients;
    }

    @Override
    public void close() {
        closeClients(clients);
        eventLoops.shutdownGracefully(0, 0, TimeUnit.MILLISECONDS).syncUninterruptibly();
    }

    private static EventLoopGroup newEventLoopGroup(Config config) {
        return switch (config.backend()) {
            case NIO -> new NioEventLoopGroup(config.eventLoops());
            case EPOLL -> {
                if (!Epoll.isAvailable()) {
                    throw new IllegalStateException("netty-bench: epoll unavailable", Epoll.unavailabilityCause());
                }
                yield new EpollEventLoopGroup(config.eventLoops());
            }
        };
    }

    private static Class<? extends Channel> clientChannelType(Config config) {
        return switch (config.backend()) {
            case NIO -> NioSocketChannel.class;
            case EPOLL -> EpollSocketChannel.class;
        };
    }

    private static void closeClients(Http2Client[] clients) {
        for (Http2Client client : clients) {
            if (client != null) {
                client.close();
            }
        }
    }
}
