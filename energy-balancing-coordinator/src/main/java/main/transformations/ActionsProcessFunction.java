package main.transformations;

import main.models.BatteryAction;
import org.apache.flink.streaming.api.functions.windowing.ProcessAllWindowFunction;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;
import org.apache.flink.util.Collector;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class ActionsProcessFunction extends ProcessAllWindowFunction<BatteryAction, List<BatteryAction>, TimeWindow> {

    @Override
    public void process(ProcessAllWindowFunction<BatteryAction, List<BatteryAction>, TimeWindow>.Context context, Iterable<BatteryAction> iterable, Collector<List<BatteryAction>> collector) throws Exception {

        System.out.println("Processing actions window at time: " + new Date().getTime());
        if (iterable == null) {
            System.out.println("Window process is null for actions");
            collector.collect(new ArrayList<>());
        }
        collector.collect((List<BatteryAction>) iterable);
    }
}
