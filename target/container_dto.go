package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (docker *Container) BuildDTO(host *Pod) (*proto.EntityDTO, error) {
	bought, _ := docker.createCommoditiesBought(host.UUID)
	sold, _ := docker.createCommoditiesSold()
	provider := builder.CreateProvider(proto.EntityDTO_CONTAINER_POD, host.UUID)

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_CONTAINER_POD, docker.UUID).
		DisplayName(docker.Name).
		Provider(provider).
		BuysCommodities(bought).
		SellsCommodities(sold).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for container(%v/%v): %v",
			pod.Name, docker.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (docker *Container) createCommoditiesBought(podId string) (*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpuComm, _ := CreateResourceCommodity(&(docker.CPU), proto.CommodityDTO_CPU)
	result = append(result, cpuComm)

	memComm, _ := CreateResourceCommodity(&(docker.Memory), proto.CommodityDTO_MEM)
	result = append(result, memComm)

	podComm, _ := CreateKeyCommodity(podId, proto.CommodityDTO_VMPM)
	result = append(result, podComm)
	return result, nil
}

func (docker *Container) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO
	cpuComm, _ := CreateResourceCommodity(&(docker.CPU), proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	memComm, _ := CreateResourceCommodity(&(docker.Memory), proto.CommodityDTO_VMEM)
	result = append(result, memComm)

	return result, nil
}
