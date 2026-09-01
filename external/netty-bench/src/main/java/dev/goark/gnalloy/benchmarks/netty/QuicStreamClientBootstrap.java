package dev.goark.gnalloy.benchmarks.netty;

import io.netty.bootstrap.Bootstrap;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.channel.ChannelOption;
import io.netty.handler.codec.quic.QuicChannel;
import io.netty.handler.codec.quic.QuicClientCodecBuilder;
import io.netty.handler.codec.quic.QuicSslContext;
import io.netty.handler.codec.quic.QuicSslContextBuilder;
import io.netty.handler.ssl.util.InsecureTrustManagerFactory;

import java.net.InetSocketAddress;
import java.util.concurrent.TimeUnit;

final class QuicStreamClientBootstrap {
    private QuicStreamClientBootstrap() {
    }

    static QuicStreamClient[] prepareClients(InetSocketAddress address, Config config, DatagramEventLoopResources resources) throws Exception {
        QuicStreamClient[] clients = new QuicStreamClient[config.connections()];
        try {
            for (int i = 0; i < clients.length; i++) {
                clients[i] = connectClient(address, config, resources, i);
            }
            return clients;
        } catch (Exception failure) {
            QuicStreamLoadGenerator.closeClients(clients);
            throw failure;
        }
    }

    private static QuicStreamClient connectClient(
            InetSocketAddress address,
            Config config,
            DatagramEventLoopResources resources,
            int clientId) throws Exception {
        Channel datagram = openDatagramChannel(resources, config);
        try {
            QuicChannel channel = QuicChannel.newBootstrap(datagram)
                    .handler(new ChannelInboundHandlerAdapter())
                    .streamHandler(new ChannelInboundHandlerAdapter())
                    .remoteAddress(address)
                    .connect()
                    .get(config.timeout().toMillis(), TimeUnit.MILLISECONDS);
            return new QuicStreamClient(datagram, channel, makePayload(config.payload(), clientId), new byte[config.payload()]);
        } catch (Exception failure) {
            datagram.close().sync();
            throw failure;
        }
    }

    private static Channel openDatagramChannel(DatagramEventLoopResources resources, Config config) throws Exception {
        ChannelHandler codec = new QuicClientCodecBuilder()
                .sslContext(clientSslContext(config))
                .maxIdleTimeout(config.timeout().toMillis(), TimeUnit.MILLISECONDS)
                .initialMaxData(16L * 1024 * 1024)
                .initialMaxStreamDataBidirectionalLocal(6L * 1024 * 1024)
                .initialMaxStreamDataBidirectionalRemote(6L * 1024 * 1024)
                .initialMaxStreamDataUnidirectional(6L * 1024 * 1024)
                .initialMaxStreamsBidirectional(QuicLimits.bidirectionalStreamLimit(config))
                .initialMaxStreamsUnidirectional(0)
                .build();
        return new Bootstrap()
                .group(resources.group())
                .channel(resources.channelType())
                .option(ChannelOption.ALLOCATOR, PooledByteBufAllocator.DEFAULT)
                .handler(codec)
                .bind(0)
                .sync()
                .channel();
    }

    private static QuicSslContext clientSslContext(Config config) throws Exception {
        return QuicSslContextBuilder
                .forClient()
                .trustManager(InsecureTrustManagerFactory.INSTANCE)
                .applicationProtocols(config.alpnProtocols().toArray(String[]::new))
                .build();
    }

    private static byte[] makePayload(int size, int clientId) {
        byte[] payload = new byte[size];
        for (int i = 0; i < payload.length; i++) {
            payload[i] = (byte) (clientId + i);
        }
        return payload;
    }
}
