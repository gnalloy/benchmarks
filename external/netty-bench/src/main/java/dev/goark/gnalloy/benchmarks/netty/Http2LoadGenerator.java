package dev.goark.gnalloy.benchmarks.netty;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.time.Duration;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.atomic.AtomicReference;

final class Http2LoadGenerator {
    private Http2LoadGenerator() {
    }

    static BenchmarkResult run(InetSocketAddress address, Config config) throws Exception {
        try (Http2ClientGroup group = Http2ClientGroup.connect(address, config)) {
            Http2Client[] clients = group.clients();
            runWarmup(clients, config);

            AtomicLong successes = new AtomicLong();
            AtomicLong errors = new AtomicLong();
            AtomicReference<Throwable> firstError = new AtomicReference<>();
            long[][] latencySamples = LatencyRecorder.samplingEnabled(config.latencySampleRate())
                    ? new long[config.connections()][]
                    : new long[0][];
            ExecutorService pool = Executors.newFixedThreadPool(config.connections());
            ResourceSnapshot resourcesBefore = ResourceSnapshot.capture();
            long started = System.nanoTime();
            for (int i = 0; i < clients.length; i++) {
                int clientId = i;
                pool.execute(() -> runMeasuredClient(
                        clients[clientId], clientId, config, successes, errors, firstError, latencySamples));
            }
            waitForWorkers(pool, config.timeout(), "HTTP/2 load", errors, firstError);
            long elapsedNanos = System.nanoTime() - started;
            long total = successes.get();
            BenchmarkResult result = result(
                    total,
                    errors.get(),
                    elapsedNanos,
                    LatencyRecorder.summarize(latencySamples),
                    firstNegotiatedProtocol(clients),
                    resourcesBefore.delta(ResourceSnapshot.capture()));
            rethrowFirstError(firstError.get());
            long expected = (long) config.connections() * config.messages();
            if (total != expected) {
                throw new IOException("netty-bench: completed " + total + " requests, want " + expected);
            }
            return result;
        }
    }

    private static void runWarmup(Http2Client[] clients, Config config) throws Exception {
        if (config.warmupMessages() <= 0) {
            return;
        }
        AtomicReference<Throwable> firstError = new AtomicReference<>();
        ExecutorService pool = Executors.newFixedThreadPool(config.connections());
        for (Http2Client client : clients) {
            pool.execute(() -> {
                try {
                    runClientMessages(client, config.warmupMessages(), config, null, new long[0]);
                } catch (Throwable failure) {
                    firstError.compareAndSet(null, failure);
                }
            });
        }
        waitForWorkers(pool, config.timeout(), "HTTP/2 warmup", null, firstError);
        rethrowFirstError(firstError.get());
    }

    private static void runMeasuredClient(
            Http2Client client,
            int clientId,
            Config config,
            AtomicLong successes,
            AtomicLong errors,
            AtomicReference<Throwable> firstError,
            long[][] latencySamples) {
        try {
            long[] samples = LatencyRecorder.newSamples(config.messages(), config.latencySampleRate());
            if (samples.length > 0) {
                latencySamples[clientId] = samples;
            }
            runClientMessages(client, config.messages(), config, successes, samples);
        } catch (Throwable failure) {
            errors.incrementAndGet();
            firstError.compareAndSet(null, failure);
        }
    }

    private static void runClientMessages(
            Http2Client client,
            int messageCount,
            Config config,
            AtomicLong successes,
            long[] latencySamples) throws IOException {
        int sampleIndex = 0;
        for (int i = 0; i < messageCount; i++) {
            boolean recordLatency = latencySamples.length > 0
                    && LatencyRecorder.shouldRecord(i, config.latencySampleRate());
            long requestStarted = recordLatency ? System.nanoTime() : 0L;
            client.exchange(config.timeout());
            if (recordLatency && sampleIndex < latencySamples.length) {
                latencySamples[sampleIndex++] = LatencyRecorder.elapsedNanos(requestStarted);
            }
            if (successes != null) {
                successes.incrementAndGet();
            }
        }
    }

    private static void waitForWorkers(
            ExecutorService pool,
            Duration timeout,
            String operation,
            AtomicLong errors,
            AtomicReference<Throwable> firstError) {
        pool.shutdown();
        try {
            if (pool.awaitTermination(timeout.toNanos(), TimeUnit.NANOSECONDS)) {
                return;
            }
            pool.shutdownNow();
            if (errors != null) {
                errors.incrementAndGet();
            }
            firstError.compareAndSet(null, new IOException("netty-bench: " + operation + " timeout"));
        } catch (InterruptedException interrupted) {
            Thread.currentThread().interrupt();
            pool.shutdownNow();
            if (errors != null) {
                errors.incrementAndGet();
            }
            firstError.compareAndSet(null, new IOException("netty-bench: " + operation + " interrupted", interrupted));
        }
    }

    private static BenchmarkResult result(
            long total,
            long errors,
            long elapsedNanos,
            LatencySummary latency,
            String negotiatedProtocol,
            ResourceDelta resources) {
        double throughput = elapsedNanos > 0 ? total * 1_000_000_000.0 / elapsedNanos : 0.0;
        double nsPerOp = total > 0 ? (double) elapsedNanos / total : 0.0;
        return new BenchmarkResult(
                total,
                errors,
                Duration.ofNanos(elapsedNanos),
                throughput,
                nsPerOp,
                negotiatedProtocol,
                latency,
                resources);
    }

    private static String firstNegotiatedProtocol(Http2Client[] clients) {
        for (Http2Client client : clients) {
            if (client != null && !client.negotiatedProtocol().isEmpty()) {
                return client.negotiatedProtocol();
            }
        }
        return "";
    }

    private static void rethrowFirstError(Throwable failure) throws Exception {
        if (failure == null) {
            return;
        }
        if (failure instanceof Exception exception) {
            throw exception;
        }
        throw new RuntimeException(failure);
    }
}
