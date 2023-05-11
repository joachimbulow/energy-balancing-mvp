package main.transformations;

import com.fasterxml.jackson.databind.ObjectMapper;
import main.models.BatteryAction;
import main.models.PemRequest;
import org.apache.flink.api.common.functions.MapFunction;


public class JsonToActionMapper implements MapFunction<String, BatteryAction> {
    ObjectMapper mapper;


    @Override
    public BatteryAction map(String s) throws Exception {
        if (mapper == null)
            mapper = new ObjectMapper();

        try {
            return mapper.readValue(s, BatteryAction.class);
        }
        catch (Exception e) {
            System.out.println("Error mapping json to batteryaction: " + e.getMessage());
        }
        return null;
    }
}