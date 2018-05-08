package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	TargetIdentifierField string = "targetIdentifier"
	Username              string = "username"
	Password              string = "password"
)

type DemoRegistrationClient struct {
	stitchingType string
}

func NewRegistrationClient(stype string) *DemoRegistrationClient {
	return &DemoRegistrationClient{
		stitchingType: stype,
	}
}

func (rClient *DemoRegistrationClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	supplyChainFactory := NewSupplyChainFactory(rClient.stitchingType)
	supplyChain, err := supplyChainFactory.createSupplyChain()
	if err != nil {
		// TODO error handling
	}
	return supplyChain
}

func (rClient *DemoRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {
	var acctDefProps []*proto.AccountDefEntry

	// target ID
	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(TargetIdentifierField, "Address",
		"IP of the target cluster master", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, targetIDAcctDefEntry)

	// username
	usernameAcctDefEntry := builder.NewAccountDefEntryBuilder(Username, "Username",
		"Username of the target cluster master", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, usernameAcctDefEntry)

	// password
	passwordAcctDefEntry := builder.NewAccountDefEntryBuilder(Password, "Password",
		"Password of the target cluster master", ".*", false, true).Create()
	acctDefProps = append(acctDefProps, passwordAcctDefEntry)

	return acctDefProps
}

func (rClient *DemoRegistrationClient) GetIdentifyingFields() string {
	return TargetIdentifierField
}

func (rClient *DemoRegistrationClient) GetActionPolicy() []*proto.ActionPolicyDTO {
	glog.V(2).Infof("Begin to build Action Policies")
	ab := builder.NewActionPolicyBuilder()
	supported := proto.ActionPolicyDTO_SUPPORTED
	recommend := proto.ActionPolicyDTO_NOT_EXECUTABLE
	notSupported := proto.ActionPolicyDTO_NOT_SUPPORTED

	//1. containerPod: move, provision; not resize;
	pod := proto.EntityDTO_CONTAINER_POD
	podCap := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	podCap[proto.ActionItemDTO_MOVE] = supported
	podCap[proto.ActionItemDTO_PROVISION] = supported
	podCap[proto.ActionItemDTO_RIGHT_SIZE] = notSupported
	addActionPolicy(ab, pod, podCap)

	//2. container: support resize; recommend provision; not move;
	container := proto.EntityDTO_CONTAINER
	containerPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	containerPolicy[proto.ActionItemDTO_RIGHT_SIZE] = supported
	//containerPolicy[proto.ActionItemDTO_RESIZE_CAPACITY] = supported
	containerPolicy[proto.ActionItemDTO_PROVISION] = recommend
	containerPolicy[proto.ActionItemDTO_MOVE] = notSupported
	addActionPolicy(ab, container, containerPolicy)

	//3. application: only recommend provision; all else are not supported
	app := proto.EntityDTO_APPLICATION
	appPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	appPolicy[proto.ActionItemDTO_PROVISION] = recommend
	appPolicy[proto.ActionItemDTO_RIGHT_SIZE] = notSupported
	appPolicy[proto.ActionItemDTO_MOVE] = notSupported
	addActionPolicy(ab, app, appPolicy)

	return ab.Create()
}

func addActionPolicy(ab *builder.ActionPolicyBuilder,
	entity proto.EntityDTO_EntityType,
	policies map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability) {

	for action, policy := range policies {
		ab.WithEntityActions(entity, action, policy)
	}
}
