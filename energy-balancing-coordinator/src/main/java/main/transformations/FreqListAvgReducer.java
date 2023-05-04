package main.transformations;

import main.models.FrequencyMeasurement;
import main.models.SystemFrequency;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.api.common.functions.ReduceFunction;

import java.time.Instant;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.Comparator;
import java.util.Date;
import java.util.List;

// Should allow for parallelization of the reduction
public class FreqListAvgReducer implements MapFunction<List<FrequencyMeasurement>, SystemFrequency> {

    @Override
    public SystemFrequency map(List<FrequencyMeasurement> frequencyMeasurements) {
        double avgFreq = frequencyMeasurements.stream().mapToDouble(m -> m.frequency).average().orElse(0.0);
        long timestamp = frequencyMeasurements.stream().map(m -> Instant.parse(m.timestamp).toEpochMilli()).max(Comparator.naturalOrder()).orElse(0L);
        return new SystemFrequency(avgFreq, timestamp);
    }
}
