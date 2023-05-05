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

    public ResponseSummary(int approvedCharge, int rejectedCharge, int approvedDischarge, int rejectedDischarge) {
            this.approvedCharge = approvedCharge;
            this.deniedCharge = rejectedCharge;
            this.approvedDischarge = approvedDischarge;
            this.deniedDischarge = rejectedDischarge;
        }
}
