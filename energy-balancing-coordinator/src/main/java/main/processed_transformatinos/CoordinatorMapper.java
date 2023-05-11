package main.processed_transformatinos;

import main.CoordinationJob;
import main.models.PemRequest;
import main.models.PemResponse;
import main.models.RequestType;
import main.models.ResponseType;
import main.redis.RedisConnectionPool;
import org.apache.flink.api.common.functions.MapFunction;
import redis.clients.jedis.Jedis;

public class CoordinatorMapper implements MapFunction<PemRequest, PemResponse> {
    private static final double NOMINAL_SYSTEM_FREQUENCY = 50;

    @Override
    public PemResponse map(PemRequest pemRequest) {
        RedisConnectionPool redisPool = RedisConnectionPool.getInstance();
        Jedis jedis = redisPool.getJedisPool().getResource();

        double currentFrequency;

        try {
            currentFrequency = Double.parseDouble(jedis.get(CoordinationJob.REDIS_FREQUENCY_KEY));
            jedis.close();
            //currentInertia = Double.parseDouble(jedis.get(CoordinationJob.REDIS_INERTIA_KEY)); REMOVED UNTIL FURTHER NOTICE
        }
        catch (Exception e) {
            System.out.println("Error parsing frequency or inertia value from Redis: " + e.getMessage());
            jedis.close();
            return null;
        }

        ResponseType responseType;
        if (currentFrequency < NOMINAL_SYSTEM_FREQUENCY) {
            responseType = pemRequest.requestType == RequestType.CHARGE ? ResponseType.GRANTED : ResponseType.DENIED;
        }
        else {
            responseType = pemRequest.requestType == RequestType.DISCHARGE ? ResponseType.GRANTED : ResponseType.DENIED;
        }
        return new PemResponse(pemRequest.id, pemRequest.batteryId, responseType, pemRequest.requestType);
    }


}
