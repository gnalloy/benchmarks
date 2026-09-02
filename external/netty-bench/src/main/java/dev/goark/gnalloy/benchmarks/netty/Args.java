package dev.goark.gnalloy.benchmarks.netty;

import java.util.HashMap;
import java.util.Map;

final class Args {
    private Args() {
    }

    static Map<String, String> parse(String[] args) {
        Map<String, String> values = new HashMap<>();
        for (int i = 0; i < args.length; i++) {
            String key = args[i];
            if (!key.startsWith("--")) {
                throw new IllegalArgumentException("netty-bench: invalid argument " + key);
            }
            if (i + 1 >= args.length) {
                throw new IllegalArgumentException("netty-bench: missing value for " + key);
            }
            values.put(key.substring(2), args[++i]);
        }
        return values;
    }

    static int intValue(Map<String, String> values, String key, int fallback) {
        String value = values.get(key);
        if (value == null || value.isBlank()) {
            return fallback;
        }
        return Integer.parseInt(value);
    }

    static boolean booleanValue(Map<String, String> values, String key, boolean fallback) {
        String value = values.get(key);
        if (value == null || value.isBlank()) {
            return fallback;
        }
        if ("true".equalsIgnoreCase(value)) {
            return true;
        }
        if ("false".equalsIgnoreCase(value)) {
            return false;
        }
        throw new IllegalArgumentException("netty-bench: invalid boolean for --" + key + ": " + value);
    }
}
