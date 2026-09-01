package dev.goark.gnalloy.benchmarks.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufAllocator;
import io.netty.handler.codec.LengthFieldBasedFrameDecoder;

final class QuicStreamFrame {
    static final int MAX_PAYLOAD_SIZE = 0xffff;

    private QuicStreamFrame() {
    }

    static LengthFieldBasedFrameDecoder decoder() {
        return new LengthFieldBasedFrameDecoder(MAX_PAYLOAD_SIZE + Short.BYTES, 0, Short.BYTES, 0, Short.BYTES);
    }

    static ByteBuf encode(ByteBufAllocator allocator, byte[] payload) {
        ByteBuf out = allocator.buffer(Short.BYTES + payload.length, Short.BYTES + payload.length);
        out.writeShort(payload.length);
        out.writeBytes(payload);
        return out;
    }

    static ByteBuf encode(ByteBufAllocator allocator, ByteBuf payload) {
        int size = payload.readableBytes();
        ByteBuf out = allocator.buffer(Short.BYTES + size, Short.BYTES + size);
        out.writeShort(size);
        out.writeBytes(payload, payload.readerIndex(), size);
        return out;
    }
}
