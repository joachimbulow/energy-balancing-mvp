package main.transformations;

import com.fasterxml.jackson.databind.ObjectMapper;
import main.models.InertiaMeasurement;
import org.apache.flink.api.common.functions.MapFunction;

public final class JsonToInertiaMapper implements MapFunction<String, InertiaMeasurement> {
    ObjectMapper mapper;

    @Override
    public InertiaMeasurement map(String s) throws Exception {
        if (mapper == null)
            mapper = new ObjectMapper();

        try {
            return mapper.readValue(s, InertiaMeasurement.class);
        }
        catch (Exception e) {
            System.out.println("Error mapping json to inertia: " + e.getMessage());
        }
        return null;
    }
}
