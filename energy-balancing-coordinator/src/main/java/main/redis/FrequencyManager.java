package main.redis;

import main.CoordinationJob;
import redis.clients.jedis.Jedis;

import java.util.Date;

public class FrequencyManager {
    private static FrequencyManager instance = null;

    private static Date lastFrequencyUpdate;
    private static double lastFrequency;

    private FrequencyManager() {
    }

    public static FrequencyManager getInstance() {
        if (instance == null) {
            instance = new FrequencyManager();
        }
        return instance;
    }

    public double getFrequency() {
        if (lastFrequencyUpdate == null || lastFrequencyUpdate.before(new Date(System.currentTimeMillis() - 1000))) {
            updateFrequency(lastFrequency);
        }

        return lastFrequency;
    }

    synchronized public void updateFrequency(double frequency) {
        // In case of multiple threads entering, only one should update the frequency
        if (lastFrequencyUpdate != null && !lastFrequencyUpdate.before(new Date(System.currentTimeMillis() - 1000))) {
            return;
        }
        Jedis jedis = RedisConnectionPool.getInstance().getJedisPool().getResource();
        lastFrequency = Double.parseDouble(jedis.get(CoordinationJob.REDIS_FREQUENCY_KEY));
        jedis.close();
    }
}
