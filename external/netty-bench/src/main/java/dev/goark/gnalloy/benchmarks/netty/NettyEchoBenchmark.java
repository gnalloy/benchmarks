package dev.goark.gnalloy.benchmarks.netty;

import java.net.InetSocketAddress;
import java.util.concurrent.CountDownLatch;

public final class NettyEchoBenchmark {
    private NettyEchoBenchmark() {
    }

    public static void main(String[] args) throws Exception {
        Config config = Config.parse(args);
        config.validate();
        if (config.serverOnly()) {
            runServerOnly(config);
            return;
        }
        BenchmarkResult result = run(config);
        if (result.totalRequests() > 0) {
            BenchmarkOutput.write(config, result);
        }
    }

    static BenchmarkResult run(Config config) throws Exception {
        config.validate();
        if (config.udpEcho()) {
            try (DatagramEchoServer server = DatagramEchoServer.start(config)) {
                return DatagramLoadGenerator.run(server.address(), config);
            }
        }
        if (config.http3Family()) {
            try (Http3Server server = Http3Server.start(config)) {
                return Http3LoadGenerator.run(server.address(), config);
            }
        }
        if (config.quicStreamFamily()) {
            try (QuicStreamServer server = QuicStreamServer.start(config)) {
                return QuicStreamLoadGenerator.run(server.address(), config);
            }
        }
        try (EchoServer server = EchoServer.start(config)) {
            InetSocketAddress address = server.address();
            if (config.http1Family()) {
                return Http1LoadGenerator.run(address, config);
            }
            if (config.http2Family()) {
                return Http2LoadGenerator.run(address, config);
            }
            return LoadGenerator.run(address, config);
        }
    }

    private static void runServerOnly(Config config) throws Exception {
        DatagramEchoServer server = DatagramEchoServer.start(config);
        Thread shutdownHook = new Thread(() -> closeServer(server), "netty-bench-shutdown");
        Runtime.getRuntime().addShutdownHook(shutdownHook);
        try {
            InetSocketAddress address = server.address();
            System.out.printf(
                    "serverReady=true framework=netty protocol=%s addr=%s:%d%n",
                    config.protocol(), address.getHostString(), address.getPort());
            System.out.flush();
            new CountDownLatch(1).await();
        } finally {
            try {
                Runtime.getRuntime().removeShutdownHook(shutdownHook);
            } catch (IllegalStateException ignored) {
            }
            server.close();
        }
    }

    private static void closeServer(DatagramEchoServer server) {
        try {
            server.close();
        } catch (InterruptedException interrupted) {
            Thread.currentThread().interrupt();
        }
    }
}
