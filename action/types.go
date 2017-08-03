package action

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
)

type TurboActionType string

const (
	ActionMove      TurboActionType = "move"
	ActionProvision TurboActionType = "provision"
	ActionUnknown   TurboActionType = "unknown"
)

type TurboExecutor interface {
	Execute(actionItem *proto.ActionItemDTO, progressTracker sdkprobe.ActionProgressTracker) (*proto.ActionResult, error)
}
