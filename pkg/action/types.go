package action

import (
	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type TurboActionType string

const (
	ActionMovePod         TurboActionType = "movePod"
	ActionMoveVM          TurboActionType = "moveVirtualMachine"
	ActionResizeContainer TurboActionType = "resizeContainer"
	ActionResizeVM        TurboActionType = "resizeVirtualMachine"
	ActionUnknown         TurboActionType = "unknown"
)

type TurboExecutor interface {
	Execute(actionItem *proto.ActionItemDTO, progressTracker sdkprobe.ActionProgressTracker) error
}
