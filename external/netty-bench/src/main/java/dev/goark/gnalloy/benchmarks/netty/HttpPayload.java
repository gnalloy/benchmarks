package dev.goark.gnalloy.benchmarks.netty;

final class HttpPayload {
    private HttpPayload() {
    }

    static byte[] body(int size) {
        byte[] body = new byte[size];
        for (int i = 0; i < body.length; i++) {
            body[i] = (byte) i;
        }
        return body;
    }
}
