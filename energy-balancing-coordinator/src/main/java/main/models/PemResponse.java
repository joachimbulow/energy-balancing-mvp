package main.models;


public class PemResponse {
    public String  id;
    public String batteryId;
    public ResponseType  type;

    public RequestType originalRequestType;

    public PemResponse(String id, String batteryId, ResponseType type, RequestType originalRequestType) {
        this.id = id;
        this.batteryId = batteryId;
        this.type = type;
        this.originalRequestType = originalRequestType;
    }
}
