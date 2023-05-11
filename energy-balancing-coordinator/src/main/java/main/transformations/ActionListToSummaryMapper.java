package main.transformations;

import main.models.*;
import org.apache.flink.api.common.functions.MapFunction;

import java.util.List;

public class ActionListToSummaryMapper implements MapFunction<List<BatteryAction>, ActionSummary> {
    @Override
    public ActionSummary map(List<BatteryAction> actions) throws Exception {
        int chargeActions = actions.stream().filter(action -> action.actionType.equals(RequestType.CHARGE)).toArray().length;
        int dischargeActions = actions.stream().filter(action -> action.actionType.equals(RequestType.DISCHARGE)).toArray().length;
        return new ActionSummary(chargeActions, dischargeActions);
    }
}
