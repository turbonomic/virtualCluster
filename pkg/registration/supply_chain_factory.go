package registration

import (
	//"fmt"

	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	cpuType     proto.CommodityDTO_CommodityType = proto.CommodityDTO_CPU
	memType     proto.CommodityDTO_CommodityType = proto.CommodityDTO_MEM
	vCpuType    proto.CommodityDTO_CommodityType = proto.CommodityDTO_VCPU
	vMemType    proto.CommodityDTO_CommodityType = proto.CommodityDTO_VMEM
	clusterType proto.CommodityDTO_CommodityType = proto.CommodityDTO_CLUSTER

	// Application Commodity is an AccessCommodity, bind the seller to the buyer
	appCommType     proto.CommodityDTO_CommodityType = proto.CommodityDTO_APPLICATION
	transactionType proto.CommodityDTO_CommodityType = proto.CommodityDTO_TRANSACTION
	vmPMAccessType  proto.CommodityDTO_CommodityType = proto.CommodityDTO_VMPM_ACCESS
	responseTimeType proto.CommodityDTO_CommodityType = proto.CommodityDTO_RESPONSE_TIME

	fakeKey string = "fake"

	CpuTemplateComm         *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &cpuType}
	MemTemplateComm         *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &memType}
	vCpuTemplateComm        *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vCpuType}
	vMemTemplateComm        *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vMemType}
	clusterTemplateComm     *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &clusterType}
	applicationTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &appCommType}
	transactionTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &transactionType}
	vmpmAccessTemplateComm  *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &vmPMAccessType}
	responseTimeTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{Key: &fakeKey, CommodityType: &responseTimeType}
)

type SupplyChainFactory struct {
}

func NewSupplyChainFactory(stype string) *SupplyChainFactory {
	return &SupplyChainFactory{}
}

func (f *SupplyChainFactory) createSupplyChain() ([]*proto.TemplateDTO, error) {
	//Physical Machine
	pmSupplyChainNode, err := f.buildPMSupply()
	if err != nil {
		return nil, err
	}

	// Virtual Machine
	vmSupplyChainNode, err := f.buildVMSupply()
	if err != nil {
		return nil, err
	}

	// Pod supply chain builder
	podSupplyChainNode, err := f.buildPodSupply()
	if err != nil {
		return nil, err
	}

	// Container
	containerSupplyChainNode, err := f.buildContainerSupply()
	if err != nil {
		return nil, err
	}

	// Application supply chain builder
	appSupplyChainNode, err := f.buildApplicationSupply()
	if err != nil {
		return nil, err
	}

	// Virtual application supply chain builder
	vAppSupplyChainNode, err := f.buildVirtualApplicationSupply()
	if err != nil {
		return nil, err
	}

	supplyChainBuilder := supplychain.NewSupplyChainBuilder()
	supplyChainBuilder.Top(vAppSupplyChainNode)
	supplyChainBuilder.Entity(appSupplyChainNode)
	supplyChainBuilder.Entity(containerSupplyChainNode)
	supplyChainBuilder.Entity(podSupplyChainNode)
	supplyChainBuilder.Entity(vmSupplyChainNode)
	supplyChainBuilder.Entity(pmSupplyChainNode)

	return supplyChainBuilder.Create()
}

func (f *SupplyChainFactory) buildPMSupply() (*proto.TemplateDTO, error) {
	nodeSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_PHYSICAL_MACHINE)
	nodeSupplyChainNodeBuilder = nodeSupplyChainNodeBuilder.
		Sells(CpuTemplateComm).
		Sells(MemTemplateComm).
		Sells(clusterTemplateComm)

	return nodeSupplyChainNodeBuilder.Create()
}

func (f *SupplyChainFactory) buildVMSupply() (*proto.TemplateDTO, error) {
	nodeSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_MACHINE)
	nodeSupplyChainNodeBuilder = nodeSupplyChainNodeBuilder.
		Sells(vCpuTemplateComm).
		Sells(vMemTemplateComm).
		Sells(clusterTemplateComm).
		Provider(proto.EntityDTO_PHYSICAL_MACHINE, proto.Provider_HOSTING).
		Buys(CpuTemplateComm).
		Buys(MemTemplateComm).
		Buys(clusterTemplateComm)
		// TODO we will re-include provisioned commodities sold by node later.
		//Sells(cpuProvisionedTemplateComm).
		//Sells(memProvisionedTemplateComm)

	return nodeSupplyChainNodeBuilder.Create()
}

func (f *SupplyChainFactory) buildPodSupply() (*proto.TemplateDTO, error) {
	// Pod supply chain node builder
	podSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_CONTAINER_POD)
	podSupplyChainNodeBuilder = podSupplyChainNodeBuilder.
		Sells(vCpuTemplateComm).
		Sells(vMemTemplateComm).
		Sells(vmpmAccessTemplateComm).
		Provider(proto.EntityDTO_VIRTUAL_MACHINE, proto.Provider_HOSTING).
		Buys(vCpuTemplateComm).
		Buys(vMemTemplateComm).
		Buys(clusterTemplateComm)

	return podSupplyChainNodeBuilder.Create()
}

func (f *SupplyChainFactory) buildContainerSupply() (*proto.TemplateDTO, error) {
	containerBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_CONTAINER)

	containerBuilder = containerBuilder.
		Sells(vCpuTemplateComm).
		Sells(vMemTemplateComm).
		Sells(applicationTemplateComm).
		Provider(proto.EntityDTO_CONTAINER_POD, proto.Provider_HOSTING).
		Buys(vmpmAccessTemplateComm).
		Buys(vCpuTemplateComm).
		Buys(vMemTemplateComm)

	return containerBuilder.Create()
}

func (f *SupplyChainFactory) buildApplicationSupply() (*proto.TemplateDTO, error) {
	// Application supply chain builder
	appSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_APPLICATION)
	appSupplyChainNodeBuilder = appSupplyChainNodeBuilder.
		Sells(transactionTemplateComm).
		Sells(responseTimeTemplateComm).
		Provider(proto.EntityDTO_CONTAINER, proto.Provider_HOSTING).
		Buys(vCpuTemplateComm).
		Buys(vMemTemplateComm).
		Buys(applicationTemplateComm)

	return appSupplyChainNodeBuilder.Create()
}

func (f *SupplyChainFactory) buildVirtualApplicationSupply() (*proto.TemplateDTO, error) {
	vAppSupplyChainNodeBuilder := supplychain.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_APPLICATION)
	vAppSupplyChainNodeBuilder = vAppSupplyChainNodeBuilder.
		Sells(transactionTemplateComm).
		Sells(responseTimeTemplateComm).
		Provider(proto.EntityDTO_APPLICATION, proto.Provider_LAYERED_OVER).
		Buys(transactionTemplateComm).
		Buys(responseTimeTemplateComm)
	return vAppSupplyChainNodeBuilder.Create()
}
