package main.models;


public class PemResponse {
    String  id;
    String batteryId;
    ResponseType  type;

    public PemResponse(String id, String batteryId, ResponseType type) {
        this.id = id;
        this.batteryId = batteryId;
        this.type = type;
    }
}
