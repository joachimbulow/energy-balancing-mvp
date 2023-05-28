package main.models;

public class ResponseSummary {
    public int approvedCharge;
    public int deniedCharge;
    public int approvedDischarge;
    public int deniedDischarge;

    public ResponseSummary() {
        this.approvedCharge = 0;
        this.deniedCharge = 0;
        this.approvedDischarge = 0;
        this.deniedDischarge = 0;
    }

    public ResponseSummary(int approvedCharge, int deniedCharge, int approvedDischarge, int deniedDischarge) {
            this.approvedCharge = approvedCharge;
            this.deniedCharge = deniedCharge;
            this.approvedDischarge = approvedDischarge;
            this.deniedDischarge = deniedDischarge;
        }
}
