package stitching

import (
	"fmt"
	"strings"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"

	"github.com/golang/glog"
)

const (
	UUID StitchingPropertyType = "UUID"
	IP   StitchingPropertyType = "IP"

	// The property used for node property and replacement entity metadata
	proxyVMIP       string = "Proxy_VM_IP"
	proxyVMUUID     string = "Proxy_VM_UUID"
	PodID           string = "POD"
	ContainerID     string = "CNT"
	ContainerFullID string = "CNTFULL"
	ContainerIDlen  int    = 12

	// The default namespace of entity property
	DefaultPropertyNamespace string = "DEFAULT"

	// The attribute used for stitching with other probes (e.g., prometurbo) with app and service
	AppStitchingAttr string = "IP"
)

// The property type that is used for stitching. For example "UUID", "IP address".
type StitchingPropertyType string

type StitchingManager struct {
	// key: node name; value: UID or IP for stitching
	nodeStitchingIDMap map[string]string

	// The property used for stitching.
	stitchType StitchingPropertyType
}

func NewStitchingManager(pType StitchingPropertyType) *StitchingManager {
	if pType != UUID && pType != IP {
		glog.Errorf("Wrong stitching type: %v, only [%v, %v] are acceptable", pType, UUID, IP)
	}

	return &StitchingManager{
		stitchType:         pType,
		nodeStitchingIDMap: make(map[string]string),
	}
}

func (s *StitchingManager) GetStitchType() StitchingPropertyType {
	return s.stitchType
}

// Get the stitching value based on given nodeName.
// Return localTestStitchingValue if it is a local testing.
func (s *StitchingManager) GetStitchingValue(nodeName string) (string, error) {
	sid, exist := s.nodeStitchingIDMap[nodeName]
	if !exist {
		err := fmt.Errorf("Failed to get stitching value for node %v, type=%v", nodeName, s.stitchType)
		glog.Error(err.Error())
		return "", err
	}

	return sid, nil
}

// Build the stitching node property for entity based on the given node name, and purpose.
//   two purposes: "stitching" and "reconcile".
//       stitching: is to stitch Pod to the real-VM;
//       reconcile: is to merge the proxy-VM to the real-VM;
func (s *StitchingManager) BuildDTOProperty(nodeName string, isForReconcile bool) (*proto.EntityDTO_EntityProperty, error) {
	propertyNamespace := DefaultPropertyNamespace
	propertyName := s.getPropertyName(isForReconcile)
	propertyValue, err := s.GetStitchingValue(nodeName)
	if err != nil {
		return nil, fmt.Errorf("Failed to build entity stitching property: %s", err)
	}
	return &proto.EntityDTO_EntityProperty{
		Namespace: &propertyNamespace,
		Name:      &propertyName,
		Value:     &propertyValue,
	}, nil
}

// Stitch one entity with a list of VMs.
func (s *StitchingManager) BuildDTOLayerOverProperty(nodeNames []string) (*proto.EntityDTO_EntityProperty, error) {
	propertyNamespace := DefaultPropertyNamespace
	propertyName := s.getStitchingPropertyName()

	values := []string{}
	for _, nodeName := range nodeNames {
		value, err := s.GetStitchingValue(nodeName)
		if err != nil {
			glog.Errorf("Failed to build DTO stitching property: %v", err)
			return nil, err
		}
		values = append(values, value)
	}
	propertyValue := strings.Join(values, ",")

	return &proto.EntityDTO_EntityProperty{
		Namespace: &propertyNamespace,
		Name:      &propertyName,
		Value:     &propertyValue,
	}, nil
}

// Get the property name based on whether it is a stitching or reconciliation.
func (s *StitchingManager) getPropertyName(isForReconcile bool) string {
	if isForReconcile {
		return s.getReconciliationPropertyName()
	}

	return s.getStitchingPropertyName()
}

// Get the name of property for entities reconciliation.
func (s *StitchingManager) getReconciliationPropertyName() string {
	if s.stitchType == UUID {
		return proxyVMUUID
	}

	return proxyVMIP
}

// Get the name of property for entities stitching.
func (s *StitchingManager) getStitchingPropertyName() string {
	if s.stitchType == UUID {
		return supplychain.SUPPLY_CHAIN_CONSTANT_UUID
	}
	return supplychain.SUPPLY_CHAIN_CONSTANT_IP_ADDRESS
}

// Create the meta data that will be used during the reconciliation process.
// This seems only applicable for VirtualMachines.
func (s *StitchingManager) GenerateReconciliationMetaData() (*proto.EntityDTO_ReplacementEntityMetaData, error) {
	replacementEntityMetaDataBuilder := builder.NewReplacementEntityMetaDataBuilder()
	switch s.stitchType {
	case UUID:
		replacementEntityMetaDataBuilder.Matching(proxyVMUUID).MatchingExternal(supplychain.VM_UUID)
	case IP:
		replacementEntityMetaDataBuilder.Matching(proxyVMIP).MatchingExternal(supplychain.VM_IP)
	default:
		return nil, fmt.Errorf("stitching property type %s is not supported", s.stitchType)
	}
	usedAndCapacityPropertyNames := []string{builder.PropertyCapacity, builder.PropertyUsed}
	vcpuUsedAndCapacityPropertyNames := []string{builder.PropertyCapacity, builder.PropertyUsed, builder.PropertyPeak}
	capacityOnlyPropertyNames := []string{builder.PropertyCapacity}
	replacementEntityMetaDataBuilder.PatchSellingWithProperty(proto.CommodityDTO_CLUSTER, capacityOnlyPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VMPM_ACCESS, capacityOnlyPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VCPU, vcpuUsedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VMEM, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VCPU_REQUEST, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VMEM_REQUEST, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VCPU_LIMIT_QUOTA, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VMEM_LIMIT_QUOTA, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VCPU_REQUEST_QUOTA, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VMEM_REQUEST_QUOTA, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_NUMBER_CONSUMERS, usedAndCapacityPropertyNames).
		PatchSellingWithProperty(proto.CommodityDTO_VSTORAGE, usedAndCapacityPropertyNames)
	meta := replacementEntityMetaDataBuilder.Build()
	return meta, nil
}
