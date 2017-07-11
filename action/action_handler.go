package action

import (
	"fmt"
	"github.com/golang/glog"
	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type ActionHandler struct {
	actionExecutors map[TurboActionType]TurboExecutor
	stop            chan struct{}
}

func NewActionHandler(stop chan struct{}) *ActionHandler {
	executors := make(map[TurboActionType]TurboExecutor)

	return &ActionHandler{
		stop:            stop,
		actionExecutors: executors,
	}
}

func (h *ActionHandler) failedResult(msg string) *proto.ActionResult {

	state := proto.ActionResponseState_FAILED
	progress := int32(0)

	res := &proto.ActionResponse{
		ActionResponseState: &state,
		Progress:            &progress,
		ResponseDescription: &msg,
	}

	return &proto.ActionResult{
		Response: res,
	}
}

func (h *ActionHandler) ExecuteAction(
	actionDTO *proto.ActionExecutionDTO,
	accountValue []*proto.AccountValue,
	progressTracker sdkprobe.ActionProgressTracker) (*proto.ActionResult, error) {

	actionItems := actionDTO.GetActionItem()
	action := actionItems[0]

	actionType, err := getActionType(action)
	if err != nil {
		msg := fmt.Sprintf("failed to get Action Type:%v", err.Error())
		glog.Error(msg)
		result := h.failedResult(msg)
		return result, nil
	}

	executor, exist := h.actionExecutors[actionType]
	if !exist {
		msg := fmt.Sprintf("action type [%v] is not supported", actionType)
		glog.Error(msg)
		result := h.failedResult(msg)
		return result, nil
	}

	result, err := executor.Execute(action)
	if err != nil {
		msg := fmt.Sprintf("Action failed: %v", err.Error())
		glog.Error(msg)
		result := h.failedResult(msg)
		return result, nil
	}

	return result, nil
}

func getActionType(action *proto.ActionItemDTO) (TurboActionType, error) {
	atype := action.GetActionType()
	switch atype {
	case proto.ActionItemDTO_MOVE:
		return ActionMove, nil
	case proto.ActionItemDTO_PROVISION:
		return ActionProvision, nil
	default:
		return ActionUnknown, fmt.Errorf("Action [%v] is not supported", atype)
	}
}
