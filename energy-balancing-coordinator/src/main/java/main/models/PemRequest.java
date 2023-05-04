package main.models;

public class PemRequest {
    public String id;
    public String batteryId;
    public RequestType requestType;

    public PemRequest() {
    }

    public PemRequest(String id, String batteryId, RequestType requestType) {
        this.id = id;
        this.batteryId = batteryId;
        this.requestType = requestType;
    }
}

