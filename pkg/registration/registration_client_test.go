package registration

import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"testing"
)

func xcheck(expected map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability,
	elements []*proto.ActionPolicyDTO_ActionPolicyElement) error {

	if len(expected) != len(elements) {
		return fmt.Errorf("length not equal: %d Vs. %d", len(expected), len(elements))
	}

	for _, e := range elements {
		action := e.GetActionType()
		capability := e.GetActionCapability()
		p, exist := expected[action]
		if !exist {
			return fmt.Errorf("action type(%v) not exist.", action)
		}

		if p != capability {
			return fmt.Errorf("action(%v) policy mismatch %v Vs %v", action, capability, p)
		}
	}

	return nil
}

func TestK8sRegistrationClient_GetActionPolicy(t *testing.T) {
	reg := NewRegClient("mock")

	supported := proto.ActionPolicyDTO_SUPPORTED
	recommend := proto.ActionPolicyDTO_NOT_EXECUTABLE
	notSupported := proto.ActionPolicyDTO_NOT_SUPPORTED

	node := proto.EntityDTO_VIRTUAL_MACHINE
	pod := proto.EntityDTO_CONTAINER_POD
	container := proto.EntityDTO_CONTAINER
	app := proto.EntityDTO_APPLICATION_COMPONENT
	service := proto.EntityDTO_SERVICE

	move := proto.ActionItemDTO_MOVE
	resize := proto.ActionItemDTO_RIGHT_SIZE
	provision := proto.ActionItemDTO_PROVISION
	suspend := proto.ActionItemDTO_SUSPEND
	scale := proto.ActionItemDTO_SCALE

	expected_pod := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expected_pod[move] = supported
	expected_pod[resize] = notSupported
	expected_pod[provision] = supported
	expected_pod[suspend] = supported

	expected_container := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expected_container[move] = notSupported
	expected_container[resize] = supported
	expected_container[provision] = recommend
	expected_container[suspend] = recommend

	expected_app := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expected_app[move] = notSupported
	expected_app[resize] = recommend
	expected_app[provision] = recommend
	expected_app[suspend] = recommend

	expected_service := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expected_service[move] = notSupported
	expected_service[resize] = notSupported
	expected_service[provision] = notSupported
	expected_service[suspend] = notSupported

	expected_node := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	expected_node[resize] = notSupported
	expected_node[provision] = supported
	expected_node[suspend] = supported
	expected_node[scale] = notSupported

	policies := reg.GetActionPolicy()

	for _, item := range policies {
		entity := item.GetEntityType()
		expected := expected_pod

		if entity == pod {
			expected = expected_pod
		} else if entity == container {
			expected = expected_container
		} else if entity == app {
			expected = expected_app
		} else if entity == node {
			expected = expected_node
		} else if entity == service {
			expected = expected_service
		} else {
			t.Errorf("Unknown entity type: %v", entity)
		}

		if err := xcheck(expected, item.GetPolicyElement()); err != nil {
			t.Errorf("Failed action policy check for entity(%v) %v", entity, err)
		}
	}
}
