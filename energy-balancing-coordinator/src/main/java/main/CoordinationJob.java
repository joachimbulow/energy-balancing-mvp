package main;

import main.models.*;
import main.processed_transformatinos.CoordinatorMapper;
import main.sinks.InfluxDBPointMapper;
import main.sinks.RedisSinkFunction;
import main.transformations.*;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.restartstrategy.RestartStrategies;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.api.java.typeutils.TypeExtractor;
import org.apache.flink.connector.base.DeliveryGuarantee;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.api.windowing.assigners.SlidingProcessingTimeWindows;

import org.apache.flink.streaming.api.windowing.time.Time;
import org.apache.flink.streaming.connectors.influxdb.InfluxDBConfig;
import org.apache.flink.streaming.connectors.influxdb.InfluxDBPoint;
import org.apache.flink.streaming.connectors.influxdb.InfluxDBSink;

import java.util.List;
import java.util.Optional;

public class CoordinationJob {
    public final static String FREQUENCY_MEASUREMENTS_TOPIC = "frequency_measurements";
    public final static String INERTIA_MEASUREMENTS_TOPIC = "inertia_measurements";
    public final static String PEM_REQUESTS_TOPIC = "pem_requests";
    public final static String PEM_RESPONSES_TOPIC = "pem_responses";

    public final static String BATTERY_ACTIONS_TOPIC = "battery_actions";
    public final static String REDIS_FREQUENCY_KEY = "frequency";
    public final static String REDIS_INERTIA_KEY = "inertia";


    public static String KAFKA_BOOTSTRAP_SERVERS = "127.0.0.1:29092";

    public static String INFLUX_URL = "http://localhost:8086";

    public static void main(String[] args) throws Exception {

        System.out.println("Starting Flink job");

        // Override with environment variables if set
        KAFKA_BOOTSTRAP_SERVERS = Optional.ofNullable(System.getenv("KAFKA_BOOTSTRAP_SERVERS")).orElse(KAFKA_BOOTSTRAP_SERVERS);
        INFLUX_URL = Optional.ofNullable(System.getenv("INFLUX_URL")).orElse(INFLUX_URL);
        System.out.println("Kafka bootstrap servers: " + KAFKA_BOOTSTRAP_SERVERS);
        System.out.println("Influx URL: " + INFLUX_URL);

        // For prod use the below
        final StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        // For testing locally
        //final StreamExecutionEnvironment env = StreamExecutionEnvironment.createLocalEnvironment();

        env.setParallelism(15);

        // Custom settings
        env.setRestartStrategy(RestartStrategies.fixedDelayRestart(5, // Gracious amount of restarts
                org.apache.flink.api.common.time.Time.seconds(15) // delay between attempts
        ));

        // Sources
        KafkaSource<String> inertiaSource = KafkaSource.<String>builder().setBootstrapServers(KAFKA_BOOTSTRAP_SERVERS).setTopics(INERTIA_MEASUREMENTS_TOPIC)
                //.setGroupId("consumerGroupid")
                .setStartingOffsets(OffsetsInitializer.latest()) //.earliest() to only read from whenever job is started
                .setValueOnlyDeserializer(new SimpleStringSchema()).build();

        KafkaSource<String> frequencySource = KafkaSource.<String>builder().setBootstrapServers(KAFKA_BOOTSTRAP_SERVERS).setTopics(FREQUENCY_MEASUREMENTS_TOPIC).setStartingOffsets(OffsetsInitializer.latest()).setValueOnlyDeserializer(new SimpleStringSchema()).build();

        KafkaSource<String> requestsSource = KafkaSource.<String>builder().setBootstrapServers(KAFKA_BOOTSTRAP_SERVERS).setTopics(PEM_REQUESTS_TOPIC).setStartingOffsets(OffsetsInitializer.latest()).setValueOnlyDeserializer(new SimpleStringSchema()).build();

        KafkaSource<String> batteryActionsSource = KafkaSource.<String>builder().setBootstrapServers(KAFKA_BOOTSTRAP_SERVERS).setTopics(BATTERY_ACTIONS_TOPIC).setStartingOffsets(OffsetsInitializer.latest()).setValueOnlyDeserializer(new SimpleStringSchema()).build();

        // Configure sinks
        InfluxDBConfig influxDbConfig = InfluxDBConfig.builder(INFLUX_URL, "admin", "admin", "influx").build();

        // ## Streams

        // # Inertia
        DataStream<String> rawInertiaStream = env.fromSource(inertiaSource, WatermarkStrategy.noWatermarks(), "Inertia source");
        DataStream<InertiaMeasurement> pojoInertiaStream = rawInertiaStream.map(new JsonToInertiaMapper()).disableChaining().name("Inertia mapper");
        pojoInertiaStream.addSink(new RedisSinkFunction<InertiaMeasurement>()).disableChaining().name("Redis inertia sink");

        // # Frequency
        DataStream<String> rawFrequencyStream = env.fromSource(frequencySource, WatermarkStrategy.noWatermarks(), "Frequency source");
        DataStream<List<FrequencyMeasurement>> pojoFreqListStream = rawFrequencyStream.map(new JsonToFreqListMapper()).disableChaining().name("FrequencyList mapper");
        DataStream<SystemFrequency> sysFreqStream = pojoFreqListStream.map(new FreqListAvgReducer()).returns(TypeExtractor.getForClass(SystemFrequency.class)).disableChaining().name("FrequencyList reducer");

        // Sink into redis
        sysFreqStream.addSink(new RedisSinkFunction<SystemFrequency>()).disableChaining().name("Redis frequency sink");

        // Map and sink into InfluxDB
        DataStream<InfluxDBPoint> influxStream = sysFreqStream.map(new InfluxDBPointMapper<SystemFrequency>()).disableChaining().name("InfluxDB frequency mapper");
        influxStream.addSink(new InfluxDBSink(influxDbConfig)).name("InfluxDB frequency sink");

        // # Requests
        DataStream<String> rawRequestStream = env.fromSource(requestsSource, WatermarkStrategy.noWatermarks(), "Requests source").disableChaining();
        DataStream<PemRequest> pojoRequestStream = rawRequestStream.map(new JsonToRequestMapper()).disableChaining().name("Request mapper");
        DataStream<PemResponse> responseStream = pojoRequestStream.map(new CoordinatorMapper()).disableChaining().name("Coordinator mapper");

        // Consume ALL responses on ALL processing instances and reduce sink into InfluxDB
        DataStream<List<PemResponse>> timedWindowResponseStream = responseStream.windowAll(SlidingProcessingTimeWindows.of(Time.seconds(5), Time.seconds(5))).process(new RequestsProcessFunction()).disableChaining().name("Requests window process function");
        DataStream<ResponseSummary> responseSummaryStream = timedWindowResponseStream.map(new ResponseListToSummaryMapper()).disableChaining().name("Response list to summary mapper");
        DataStream<InfluxDBPoint> influxResponseStream = responseSummaryStream.map(new InfluxDBPointMapper<ResponseSummary>()).disableChaining().name("InfluxDB response summary mapper");
        influxResponseStream.addSink(new InfluxDBSink(influxDbConfig)).disableChaining().name("InfluxDB response summary sink");

        // Sink into Kafka
        DataStream<String> jsonResponseStream = responseStream.map(new PojoToJsonMapper<PemResponse>()).disableChaining().name("Response to JSON mapper");

        KafkaSink<String> kafkaSink = KafkaSink.<String>builder().setBootstrapServers(KAFKA_BOOTSTRAP_SERVERS).setRecordSerializer(KafkaRecordSerializationSchema.builder()
                .setTopic(PEM_RESPONSES_TOPIC).setValueSerializationSchema(new SimpleStringSchema()).build()).setDeliveryGuarantee(DeliveryGuarantee.NONE).build();

        jsonResponseStream.sinkTo(kafkaSink).name("Kafka response sink").disableChaining();

        // # Actions
        DataStream<String> rawActionsStream = env.fromSource(batteryActionsSource, WatermarkStrategy.noWatermarks(), "Actions source");
        DataStream<BatteryAction> pojoActionStream = rawActionsStream.map(new JsonToActionMapper()).disableChaining().name("Action mapper");

        // Sink actions into InfluxDB as well
        DataStream<List<BatteryAction>> timedWindowActionStream = pojoActionStream.windowAll(SlidingProcessingTimeWindows.of(Time.seconds(5), Time.seconds(5))).process(new ActionsProcessFunction()).disableChaining().name("Actions window process function");
        DataStream<ActionSummary> actionSummaryStream = timedWindowActionStream.map(new ActionListToSummaryMapper()).disableChaining().name("Action list to summary mapper");
        DataStream<InfluxDBPoint> influxActionStream = actionSummaryStream.map(new InfluxDBPointMapper<ActionSummary>()).disableChaining().name("InfluxDB actions summary mapper");
        influxActionStream.addSink(new InfluxDBSink(influxDbConfig)).name("InfluxDB actions summary sink");

        // Execute program, beginning computation.
        env.execute("Flink coordinator job");
    }
}
