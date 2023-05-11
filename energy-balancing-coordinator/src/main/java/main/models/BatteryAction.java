package main.models;

public class BatteryAction {
    public String id;
    public String batteryId;
    public RequestType actionType;

    public BatteryAction() {
    }

    public BatteryAction(String id, String batteryId, RequestType actionType) {
        this.id = id;
        this.batteryId = batteryId;
        this.actionType = actionType;
    }
}
