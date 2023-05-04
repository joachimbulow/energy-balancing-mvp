package main.models;

public class SystemFrequency {
    private final double frequency;
    private final long timestamp;

    public SystemFrequency(double frequency, long timestamp) {
        this.frequency = frequency;
        this.timestamp = timestamp;
    }

    public double getFrequency() {
        return frequency;
    }

    public long getTimestamp() {
        return timestamp;
    }
}
