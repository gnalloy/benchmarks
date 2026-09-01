package dev.goark.gnalloy.benchmarks.netty;

final class QuicLimits {
    private static final long MIN_BIDIRECTIONAL_STREAMS = 1024;
    private static final long STREAM_CREDIT_MARGIN = 16;

    private QuicLimits() {
    }

    static long bidirectionalStreamLimit(Config config) {
        long requestStreamsPerConnection = (long) config.messages() + config.warmupMessages();
        long workloadCapacity = requestStreamsPerConnection + STREAM_CREDIT_MARGIN;
        long connectionCapacity = (long) config.connections() * 4L;
        return Math.max(MIN_BIDIRECTIONAL_STREAMS, Math.max(workloadCapacity, connectionCapacity));
    }
}
