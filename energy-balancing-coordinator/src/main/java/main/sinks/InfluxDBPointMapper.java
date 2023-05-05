package main.sinks;

import main.models.ResponseSummary;
import main.models.SystemFrequency;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.streaming.connectors.influxdb.InfluxDBPoint;

import java.util.Date;
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
        if (in instanceof ResponseSummary) {
            Map<String, Object> fields = new HashMap<>();
            fields.put("approvedCharge", ((ResponseSummary) in).approvedCharge);
            fields.put("approvedDischarge", ((ResponseSummary) in).approvedDischarge);
            fields.put("deniedCharge", ((ResponseSummary) in).deniedCharge);
            fields.put("deniedDischarge", ((ResponseSummary) in).deniedDischarge);
            long timestamp = new Date().getTime();
            return new InfluxDBPoint("responseSummary", timestamp, null, fields);
        }
        System.out.println("Error mapping to InfluxDBPoint");
        return null;
    }
}
