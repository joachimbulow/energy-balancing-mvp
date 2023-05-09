package main.redis;
import main.CoordinationJob;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;

import java.io.IOException;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.UnknownHostException;

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
        String redisHost = CoordinationJob.REDIS_BROKER;
        int redisPort = CoordinationJob.REDIS_PORT;
        String redisPassword = null;
        System.out.println("REDIS SINK: Connecting to Redis at " + redisHost + ":" + redisPort);
        pingRedisHost(redisHost);

        jedisPool = new JedisPool(poolConfig, redisHost, redisPort, 20000, redisPassword);
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

    private boolean pingRedisHost(String host) {
        try {
            InetAddress inetAddress = InetAddress.getByName(host);
            System.out.println("REDIS SINK: " + host + " is reachable");
            return true;
        } catch (UnknownHostException e) {
            System.out.println("REDIS SINK: " + host + " is not reachable");
            System.out.println("REDIS SINK: Error while pinging Redis host: " + e.getMessage());
            return false;
        }
    }


}

