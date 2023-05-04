package main.models;

public class FrequencyMeasurement {
    public String timestamp;
    public String location;
    public double voltage;
    public double current;
    public int frequency;
    public double consumption;
    public double production;

    public FrequencyMeasurement() {
    }

    public FrequencyMeasurement(String timestamp, String location, double voltage, double current, int frequency, double consumption, double production) {
        this.timestamp = timestamp;
        this.location = location;
        this.voltage = voltage;
        this.current = current;
        this.frequency = frequency;
        this.consumption = consumption;
        this.production = production;
    }


}

