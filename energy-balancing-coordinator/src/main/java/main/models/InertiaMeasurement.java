package main.models;

import com.fasterxml.jackson.annotation.JsonProperty;

public class InertiaMeasurement {
    @JsonProperty("HourUTC")
    public String hourUTC;
    @JsonProperty("HourDK")
    public String hourDK;
    @JsonProperty("InertiaNordicGWs")
    public double inertiaNordicGWs;
    @JsonProperty("InertiaDK2GWs")
    public double inertiaDK2GWs;
    @JsonProperty("InertiaNOGWs")
    public double inertiaNOGWs;
    @JsonProperty("InertiaSEGWs")
    public double inertiaSEGWs;
    @JsonProperty("InertiaFIGWs")
    public double inertiaFIGWs;

    public InertiaMeasurement() {
    }

    public InertiaMeasurement(String hourUTC, String hourDK, double inertiaNordicGWs, double inertiaDK2GWs, double inertiaNOGWs, double inertiaSEGWs, double inertiaFIGWs) {
        this.hourUTC = hourUTC;
        this.hourDK = hourDK;
        this.inertiaNordicGWs = inertiaNordicGWs;
        this.inertiaDK2GWs = inertiaDK2GWs;
        this.inertiaNOGWs = inertiaNOGWs;
        this.inertiaSEGWs = inertiaSEGWs;
        this.inertiaFIGWs = inertiaFIGWs;
    }
}