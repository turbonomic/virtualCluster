package action

import (
	"fmt"
	"github.com/golang/glog"
	"time"

	"github.com/songbinliu/containerChain/pkg/action/executor"
	"github.com/songbinliu/containerChain/pkg/target"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type ActionHandler struct {
	cluster         *target.ClusterHandler
	actionExecutors map[TurboActionType]TurboExecutor
	stop            chan struct{}
}

func NewActionHandler(h *target.ClusterHandler, stop chan struct{}) *ActionHandler {
	executors := make(map[TurboActionType]TurboExecutor)

	handler := &ActionHandler{
		cluster:         h,
		stop:            stop,
		actionExecutors: executors,
	}

	handler.registerExecutors()

	return handler
}

func (h *ActionHandler) String() string {
	cinfo := fmt.Sprintf("%s", h.cluster.String())

	atypes := []TurboActionType{}
	for k := range h.actionExecutors {
		atypes = append(atypes, k)
	}

	return fmt.Sprintf("%s\n%v", cinfo, atypes)
}

func (h *ActionHandler) registerExecutors() {
	podMover := executor.NewPodMover(h.cluster)
	h.actionExecutors[ActionMovePod] = podMover

	containerResizer := executor.NewContainerResizer(h.cluster)
	h.actionExecutors[ActionResizeContainer] = containerResizer

	vmMover := executor.NewVirtualMachineMover(h.cluster)
	h.actionExecutors[ActionMoveVM] = vmMover

	vmResizer := executor.NewVirtualMachineMover(h.cluster)
	h.actionExecutors[ActionResizeVM] = vmResizer
}

func (h *ActionHandler) goodResult(msg string) *proto.ActionResult {

	state := proto.ActionResponseState_SUCCEEDED
	progress := int32(100)

	res := &proto.ActionResponse{
		ActionResponseState: &state,
		Progress:            &progress,
		ResponseDescription: &msg,
	}

	return &proto.ActionResult{
		Response: res,
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

	// here progressTracker is used to keep alive; executor won't really use it.
	stop := make(chan struct{})
	defer close(stop)
	go keepAlive(progressTracker, stop)

	err = executor.Execute(action, progressTracker)
	if err != nil {
		msg := fmt.Sprintf("Action failed: %v", err.Error())
		glog.Error(msg)
		result := h.failedResult(msg)
		return result, nil
	}

	result := h.goodResult("Success")
	return result, nil
}

func getActionType(action *proto.ActionItemDTO) (TurboActionType, error) {
	atype := action.GetActionType()
	object := action.GetTargetSE()
	if object == nil {
		return ActionUnknown, fmt.Errorf("TargetSE is empty.")
	}
	objectType := object.GetEntityType()

	glog.V(2).Infof("action [%v-%v] is received.", atype, objectType)

	switch atype {
	case proto.ActionItemDTO_MOVE:
		// only support move Pod, and Virtual Machine
		switch objectType {
		case proto.EntityDTO_CONTAINER_POD:
			return ActionMovePod, nil
		case proto.EntityDTO_VIRTUAL_MACHINE:
			return ActionMoveVM, nil
		}
	case proto.ActionItemDTO_RESIZE:
		switch objectType {
		case proto.EntityDTO_CONTAINER:
			return ActionResizeContainer, nil
		case proto.EntityDTO_VIRTUAL_MACHINE:
			return ActionResizeVM, nil
		}
		//case proto.ActionItemDTO_PROVISION:
		//	return ActionProvision, nil
	}

	err := fmt.Errorf("Action [%v-%v] is not supported.", atype, objectType)
	return ActionUnknown, err
}

func keepAlive(tracker sdkprobe.ActionProgressTracker, stop chan struct{}) {
	//TODO: add timeout
	go func() {
		var progress int32 = 0
		state := proto.ActionResponseState_IN_PROGRESS

		for {
			progress = progress + 1
			if progress > 99 {
				progress = 99
			}

			tracker.UpdateProgress(state, "in progress", progress)

			t := time.NewTimer(time.Second * 3)
			select {
			case <-stop:
				return
			case <-t.C:
			}
		}
		glog.V(3).Infof("action keepAlive goroutine exit.")
	}()
}
