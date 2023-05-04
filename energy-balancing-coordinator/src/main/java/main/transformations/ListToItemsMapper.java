package main.transformations;

import org.apache.flink.api.common.functions.FlatMapFunction;
import org.apache.flink.util.Collector;

import java.util.List;

public final class ListToItemsMapper<T> implements FlatMapFunction<List<T>, T> {

    @Override
    public void flatMap(List<T> list, Collector<T> out) {
        for (T item : list) {
            out.collect(item);
        }
    }
}
