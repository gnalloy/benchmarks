package dev.goark.gnalloy.benchmarks.netty;

import io.netty.handler.codec.LengthFieldBasedFrameDecoder;

final class QuicStreamFrame {
    static final int MAX_PAYLOAD_SIZE = 0xffff;

    private QuicStreamFrame() {
    }

    static LengthFieldBasedFrameDecoder decoder() {
        return new LengthFieldBasedFrameDecoder(MAX_PAYLOAD_SIZE + Short.BYTES, 0, Short.BYTES, 0, Short.BYTES);
    }
}
