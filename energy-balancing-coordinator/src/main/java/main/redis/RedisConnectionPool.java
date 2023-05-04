package main.redis;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;

public class RedisConnectionPool {

    private static RedisConnectionPool instance = null;

    private final JedisPool jedisPool;

    private RedisConnectionPool() {
        JedisPoolConfig poolConfig = new JedisPoolConfig();
        poolConfig.setMaxTotal(10);
        poolConfig.setMaxIdle(5);
        poolConfig.setMinIdle(1); // For minimum latency ;) haha
        poolConfig.setTestOnBorrow(true);
        poolConfig.setTestOnReturn(true);

        // configure the JedisPool with your Redis instance information
        String redisHost = "localhost";
        int redisPort = 6379;
        String redisPassword = null;
        jedisPool = new JedisPool(poolConfig, redisHost, redisPort, 2000, redisPassword);
    }

    public static RedisConnectionPool getInstance() {
        if (instance == null) {
            synchronized (RedisConnectionPool.class) {
                if (instance == null) {
                    instance = new RedisConnectionPool();
                }
            }
        }
        return instance;
    }

    public JedisPool getJedisPool() {
        return jedisPool;
    }

}

