package main.transformations;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.flink.api.common.functions.MapFunction;

public class PojoToJsonMapper<T> implements MapFunction<T, String> {
    private final ObjectMapper mapper = new ObjectMapper();

    @Override
    public String map(T pojo) throws JsonProcessingException {
        try {
            return mapper.writeValueAsString(pojo);
        }
        catch (Exception e) {
            System.out.println("Error mapping pojo to json: " + e.getMessage());
        }
        return null;
    }
}
