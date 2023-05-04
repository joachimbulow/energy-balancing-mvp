package main.transformations;

import main.models.PemRequest;
import main.models.PemResponse;
import org.apache.flink.streaming.api.functions.windowing.ProcessAllWindowFunction;
import org.apache.flink.streaming.api.functions.windowing.ProcessWindowFunction;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;
import org.apache.flink.streaming.api.windowing.windows.Window;
import org.apache.flink.util.Collector;
import scala.collection.immutable.List$;

import java.util.Date;
import java.util.List;

public class RequestsProcessFunction extends ProcessAllWindowFunction<PemResponse, List<PemResponse>, TimeWindow> {

    @Override
    public void process(ProcessAllWindowFunction<PemResponse, List<PemResponse>, TimeWindow>.Context context, Iterable<PemResponse> iterable, Collector<List<PemResponse>> collector) throws Exception {

        System.out.println("Processing request window at time: " + new Date().getTime());
        collector.collect((List<PemResponse>) iterable);
    }
}
