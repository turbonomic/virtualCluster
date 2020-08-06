package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/virtualCluster/pkg/discovery/stitching"
)

const (
	TargetIdentifierField string = "targetIdentifier"
	Username              string = "username"
	Password              string = "password"
)

type DemoRegClient struct {
	stitchingType stitching.StitchingPropertyType
}

func NewRegClient(pType stitching.StitchingPropertyType) *DemoRegClient {
	return &DemoRegClient{
		stitchingType: pType,
	}
}

func (rClient *DemoRegClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	supplyChainFactory := NewSupplyChainFactory(rClient.stitchingType)
	supplyChain, err := supplyChainFactory.createSupplyChain()
	if err != nil {
		glog.Errorf("Failed to create supply chain: %v", err)
		// TODO error handling
	}
	return supplyChain
}

func (rClient *DemoRegClient) GetAccountDefinition() []*proto.AccountDefEntry {
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

func (rClient *DemoRegClient) GetIdentifyingFields() string {
	return TargetIdentifierField
}

func (rClient *DemoRegClient) GetActionPolicy() []*proto.ActionPolicyDTO {
	glog.V(3).Infof("Begin to build Action Policies")
	ab := builder.NewActionPolicyBuilder()
	supported := proto.ActionPolicyDTO_SUPPORTED
	recommend := proto.ActionPolicyDTO_NOT_EXECUTABLE
	notSupported := proto.ActionPolicyDTO_NOT_SUPPORTED

	// 1. containerPod: support move, provision and suspend; not resize;
	pod := proto.EntityDTO_CONTAINER_POD
	podPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	podPolicy[proto.ActionItemDTO_MOVE] = supported
	podPolicy[proto.ActionItemDTO_PROVISION] = supported
	podPolicy[proto.ActionItemDTO_RIGHT_SIZE] = notSupported
	podPolicy[proto.ActionItemDTO_SUSPEND] = supported

	rClient.addActionPolicy(ab, pod, podPolicy)

	// 2. container: support resize; recommend provision and suspend; not move;
	container := proto.EntityDTO_CONTAINER
	containerPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	containerPolicy[proto.ActionItemDTO_RIGHT_SIZE] = supported
	containerPolicy[proto.ActionItemDTO_PROVISION] = recommend
	containerPolicy[proto.ActionItemDTO_MOVE] = notSupported
	containerPolicy[proto.ActionItemDTO_SUSPEND] = recommend

	rClient.addActionPolicy(ab, container, containerPolicy)

	// 3. application: only recommend provision and suspend; all else are not supported
	app := proto.EntityDTO_APPLICATION_COMPONENT
	appPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	appPolicy[proto.ActionItemDTO_PROVISION] = recommend
	appPolicy[proto.ActionItemDTO_RIGHT_SIZE] = recommend
	appPolicy[proto.ActionItemDTO_MOVE] = notSupported
	appPolicy[proto.ActionItemDTO_SUSPEND] = recommend

	rClient.addActionPolicy(ab, app, appPolicy)

	// 4. service: no actions are supported
	service := proto.EntityDTO_SERVICE
	servicePolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	servicePolicy[proto.ActionItemDTO_PROVISION] = notSupported
	servicePolicy[proto.ActionItemDTO_RIGHT_SIZE] = notSupported
	servicePolicy[proto.ActionItemDTO_MOVE] = notSupported
	servicePolicy[proto.ActionItemDTO_SUSPEND] = notSupported

	rClient.addActionPolicy(ab, service, servicePolicy)

	// 5. node: support provision and suspend; not resize; do not set move
	vnode := proto.EntityDTO_VIRTUAL_MACHINE
	vnodePolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	vnodePolicy[proto.ActionItemDTO_PROVISION] = supported
	vnodePolicy[proto.ActionItemDTO_RIGHT_SIZE] = notSupported
	vnodePolicy[proto.ActionItemDTO_SCALE] = notSupported
	vnodePolicy[proto.ActionItemDTO_SUSPEND] = supported

	rClient.addActionPolicy(ab, vnode, vnodePolicy)

	return ab.Create()
}

func (rClient *DemoRegClient) addActionPolicy(ab *builder.ActionPolicyBuilder,
	entity proto.EntityDTO_EntityType,
	policies map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability) {

	for action, policy := range policies {
		ab.WithEntityActions(entity, action, policy)
	}
}

func (rclient *DemoRegClient) GetEntityMetadata() []*proto.EntityIdentityMetadata {
	glog.V(3).Infof("Begin to build EntityIdentityMetadata")

	var result []*proto.EntityIdentityMetadata

	entities := []proto.EntityDTO_EntityType{
		proto.EntityDTO_SWITCH,
		proto.EntityDTO_PHYSICAL_MACHINE,
		proto.EntityDTO_NAMESPACE,
		proto.EntityDTO_WORKLOAD_CONTROLLER,
		proto.EntityDTO_VIRTUAL_MACHINE,
		proto.EntityDTO_CONTAINER_SPEC,
		proto.EntityDTO_CONTAINER_POD,
		proto.EntityDTO_CONTAINER,
		proto.EntityDTO_APPLICATION_COMPONENT,
		proto.EntityDTO_SERVICE,
	}

	for _, etype := range entities {
		meta := rclient.newIdMetaData(etype, []string{"id"})
		result = append(result, meta)
	}

	glog.V(4).Infof("EntityIdentityMetaData: %++v", result)

	return result
}

func (rclient *DemoRegClient) newIdMetaData(etype proto.EntityDTO_EntityType, names []string) *proto.EntityIdentityMetadata {
	var data []*proto.EntityIdentityMetadata_PropertyMetadata
	for _, name := range names {
		dat := &proto.EntityIdentityMetadata_PropertyMetadata{
			Name: &name,
		}
		data = append(data, dat)
	}

	result := &proto.EntityIdentityMetadata{
		EntityType:            &etype,
		NonVolatileProperties: data,
	}

	return result
}
