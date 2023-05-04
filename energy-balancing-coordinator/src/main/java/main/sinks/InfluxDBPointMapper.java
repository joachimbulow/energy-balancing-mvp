package main.sinks;

import main.models.SystemFrequency;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.streaming.connectors.influxdb.InfluxDBPoint;

import java.util.HashMap;
import java.util.Map;

public class InfluxDBPointMapper<T> implements MapFunction<T, InfluxDBPoint> {

    @Override
    public InfluxDBPoint map(T in) {
        if (in instanceof SystemFrequency) {
            Map<String, Object> fields = new HashMap<>();
            fields.put("frequency", ((SystemFrequency) in).getFrequency());
            return new InfluxDBPoint("frequency", ((SystemFrequency) in).getTimestamp(), null, fields);
        }
        return null;
    }
}
