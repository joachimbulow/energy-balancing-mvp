package main.models;


public class PemResponse {
    public String  id;
    public String batteryId;
    public ResponseType responseType;

    public RequestType originalRequestType;

    public PemResponse(String id, String batteryId, ResponseType responseType, RequestType originalRequestType) {
        this.id = id;
        this.batteryId = batteryId;
        this.responseType = responseType;
        this.originalRequestType = originalRequestType;
    }
}
