package action

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type TurboActionType string

const (
	ActionMove      TurboActionType = "move"
	ActionProvision TurboActionType = "provision"
	ActionUnknown   TurboActionType = "unknown"
)

type TurboExecutor interface {
	Execute(actionItem *proto.ActionItemDTO) (*proto.ActionResult, error)
}
