package main.transformations;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import main.models.FrequencyMeasurement;
import org.apache.flink.api.common.functions.MapFunction;

import java.util.List;

public class JsonToFreqListMapper implements MapFunction<String, List<FrequencyMeasurement>> {
    ObjectMapper mapper;


    @Override
    public List<FrequencyMeasurement> map(String s) throws Exception {
        if (mapper == null)
            mapper = new ObjectMapper();

        try {
            return mapper.readValue(s, new TypeReference<List<FrequencyMeasurement>>(){});
        }
        catch (Exception e) {
            System.out.println("Error mapping json to List<Frequency>: " + e.getMessage());
        }
        return null;
    }
}
