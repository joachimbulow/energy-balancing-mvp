package main.redis;

import main.CoordinationJob;
import redis.clients.jedis.Jedis;

import java.util.Date;

import java.util.concurrent.locks.ReentrantLock;

public class FrequencyManager {
    private static FrequencyManager instance = null;

    private static long lastFrequencyUpdateTime;
    private static double lastFrequency;
    private final ReentrantLock lock = new ReentrantLock();

    private FrequencyManager() {
        lastFrequencyUpdateTime = 0;
        lastFrequency = 50;
    }

    public static FrequencyManager getInstance() {
        if (instance == null) {
            instance = new FrequencyManager();
        }
        return instance;
    }

    public double getFrequency() {
        if (isStale()) {
            updateFrequency();
        }

        return lastFrequency;
    }

    private boolean isStale() {
        return lastFrequencyUpdateTime == 0 || lastFrequencyUpdateTime < System.currentTimeMillis() - 1000;
    }

    public void updateFrequency() {
        // Try to acquire the lock
        if (lock.tryLock()) {
            try {
                if (!isStale()) {
                    return;
                }
                Jedis jedis = RedisConnectionPool.getInstance().getJedisPool().getResource();
                lastFrequency = Double.parseDouble(jedis.get(CoordinationJob.REDIS_FREQUENCY_KEY));
                lastFrequencyUpdateTime = System.currentTimeMillis();
                jedis.close();
            } finally {
                // Always release the lock in the final block to ensure it's released even if an exception occurs
                lock.unlock();
            }
        }
    }
}

