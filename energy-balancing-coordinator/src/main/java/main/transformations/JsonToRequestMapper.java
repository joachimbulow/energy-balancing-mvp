package main.transformations;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import main.models.FrequencyMeasurement;
import main.models.PemRequest;
import org.apache.flink.api.common.functions.MapFunction;

import java.util.List;

public class JsonToRequestMapper implements MapFunction<String, PemRequest> {
    ObjectMapper mapper;


    @Override
    public PemRequest map(String s) throws Exception {
        if (mapper == null)
            mapper = new ObjectMapper();

        try {
            return mapper.readValue(s.toString(), PemRequest.class);
        }
        catch (Exception e) {
            System.out.println("Error mapping json to pemrequest: " + e.getMessage());
        }
        return null;
    }
}