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
import io.netty.handler.codec.http.HttpClientCodec;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslHandler;

import java.net.InetSocketAddress;
import java.util.concurrent.TimeUnit;

final class Http1ClientGroup implements AutoCloseable {
    private final EventLoopGroup eventLoops;
    private final Http1Client[] clients;

    private Http1ClientGroup(EventLoopGroup eventLoops, Http1Client[] clients) {
        this.eventLoops = eventLoops;
        this.clients = clients;
    }

    static Http1ClientGroup connect(InetSocketAddress address, Config config) throws Exception {
        EventLoopGroup eventLoops = newEventLoopGroup(config);
        Http1Client[] clients = new Http1Client[config.connections()];
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
                            channel.pipeline().addLast(new HttpClientCodec());
                            channel.pipeline().addLast(new HttpObjectAggregator(maxResponseBytes(config.payload())));
                            channel.pipeline().addLast(new Http1ResponseHandler(config.payload()));
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
                clients[i] = new Http1Client(
                        channel,
                        channel.pipeline().get(Http1ResponseHandler.class),
                        config.host(),
                        negotiatedProtocol);
            }
            return new Http1ClientGroup(eventLoops, clients);
        } catch (Throwable failure) {
            closeClients(clients);
            eventLoops.shutdownGracefully(0, 0, TimeUnit.MILLISECONDS).syncUninterruptibly();
            if (failure instanceof Exception exception) {
                throw exception;
            }
            throw new RuntimeException(failure);
        }
    }

    Http1Client[] clients() {
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

    private static int maxResponseBytes(int payload) {
        return Math.addExact(payload, 16 * 1024);
    }

    private static void closeClients(Http1Client[] clients) {
        for (Http1Client client : clients) {
            if (client != null) {
                client.close();
            }
        }
    }
}
