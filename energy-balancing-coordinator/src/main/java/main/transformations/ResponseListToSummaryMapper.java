package main.transformations;

import main.models.PemResponse;
import main.models.RequestType;
import main.models.ResponseSummary;
import main.models.ResponseType;
import org.apache.flink.api.common.functions.MapFunction;

import java.util.List;

public class ResponseListToSummaryMapper implements MapFunction<List<PemResponse>, ResponseSummary> {
    @Override
    public ResponseSummary map(List<PemResponse> pemResponses) throws Exception {
        int approvedCharge = pemResponses.stream().filter(pr -> pr.originalRequestType.equals(RequestType.CHARGE) && pr.responseType.equals(ResponseType.GRANTED)).toArray().length;
        int approvedDischarge = pemResponses.stream().filter(pr -> pr.originalRequestType.equals(RequestType.DISCHARGE) && pr.responseType.equals(ResponseType.GRANTED)).toArray().length;
        int deniedCharge = pemResponses.stream().filter(pr -> pr.originalRequestType.equals(RequestType.CHARGE) && pr.responseType.equals(ResponseType.DENIED)).toArray().length;
        int deniedDischarge = pemResponses.stream().filter(pr -> pr.originalRequestType.equals(RequestType.DISCHARGE) && pr.responseType.equals(ResponseType.DENIED)).toArray().length;
        return new ResponseSummary(approvedCharge, approvedDischarge, deniedCharge, deniedDischarge);
    }
}
