package dev.goark.gnalloy.benchmarks.netty;

import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;

@ChannelHandler.Sharable
final class Http2InboundStreamHandler extends ChannelInboundHandlerAdapter {
    static final Http2InboundStreamHandler INSTANCE = new Http2InboundStreamHandler();

    private Http2InboundStreamHandler() {
    }

    @Override
    public void channelActive(ChannelHandlerContext context) {
        context.close();
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext context, Throwable cause) {
        context.close();
    }
}
