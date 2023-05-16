package main.sinks;

import main.CoordinationJob;
import main.models.InertiaMeasurement;
import main.models.SystemFrequency;
import main.redis.RedisConnectionPool;
import org.apache.flink.streaming.api.functions.ProcessFunction;
import org.apache.flink.streaming.api.functions.sink.SinkFunction;
import org.apache.flink.util.Collector;
import redis.clients.jedis.Jedis;

public class RedisSinkFunction<T> implements SinkFunction<T> {

    @Override
    public void invoke(T in) throws Exception {
        RedisConnectionPool redisPool = RedisConnectionPool.getInstance();
        Jedis jedis = redisPool.getJedisPool().getResource();

        if (in instanceof InertiaMeasurement) {
            String dk2Inertia = String.valueOf(((InertiaMeasurement) in).inertiaDK2GWs);
            jedis.set(CoordinationJob.REDIS_INERTIA_KEY, dk2Inertia);
            jedis.close();
            return;
        }

        if (in instanceof SystemFrequency) {
            String systemFrequency = String.valueOf(((SystemFrequency) in).getFrequency());
            System.out.println("Writing system frequency to redis: " + systemFrequency);
            jedis.set(CoordinationJob.REDIS_FREQUENCY_KEY, systemFrequency);
            jedis.close();
            return;
        }

        System.out.println("Error sinking to Redis - unknown type");
    }
}
