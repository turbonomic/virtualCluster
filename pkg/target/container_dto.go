package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (d *Container) Clone(newName, newId string) *Container {
	result := NewContainer(newName, newId)

	result.Memory = d.Memory
	result.CPU = d.CPU
	result.ReqMemory = d.ReqMemory
	result.ReqCPU = d.ReqCPU
	result.QPS = d.QPS

	//not copy the APP
	result.App = nil
	return result
}

// this should be called after Pods have been built
func (d *Container) GenerateApp() error {
	appName := fmt.Sprintf("app-%s", d.Name)
	appId := fmt.Sprintf("app-%s", d.UUID)
	app := NewApplication(appName, appId)

	app.CPU = d.CPU
	app.Memory = d.Memory
	app.ProviderID = d.UUID
	app.QPS = d.QPS

	d.App = app

	glog.V(4).Infof("App: %+v", app)
	return nil
}

func (d *Container) BuildAppDTO(pod *Pod) (*proto.EntityDTO, error) {
	if d.App == nil {
		return nil, fmt.Errorf("container[%s] App is nil.", d.Name)
		//d.GenerateApp()
	}

	return d.App.BuildDTO(d, pod)
}

func (docker *Container) BuildDTO(pod *Pod) (*proto.EntityDTO, error) {
	bought, _ := docker.createCommoditiesBought(pod.UUID)
	sold, _ := docker.createCommoditiesSold()
	provider := builder.CreateProvider(proto.EntityDTO_CONTAINER_POD, pod.UUID)

	entity, err := builder.
		NewEntityDTOBuilder(proto.EntityDTO_CONTAINER, docker.UUID).
		DisplayName(docker.Name).
		Provider(provider).
		BuysCommodities(bought).
		SellsCommodities(sold).
		WithPowerState(proto.EntityDTO_POWERED_ON).
		Create()

	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for container(%v/%v): %v",
			pod.Name, docker.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (docker *Container) createCommoditiesBought(podId string) ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	cpuComm, _ := CreateCommodityBoughtWithReservation(&(docker.CPU), docker.ReqCPU, proto.CommodityDTO_VCPU)
	result = append(result, cpuComm)

	memComm, _ := CreateCommodityBoughtWithReservation(&(docker.Memory), docker.ReqMemory, proto.CommodityDTO_VMEM)
	result = append(result, memComm)

	podComm, _ := CreateKeyCommodity(podId, proto.CommodityDTO_VMPM_ACCESS)
	result = append(result, podComm)
	return result, nil
}

func (docker *Container) createCommoditiesSold() ([]*proto.CommodityDTO, error) {

	var result []*proto.CommodityDTO

	resizeable := true
	cpuComm, _ := CreateResourceCommodityResize(&(docker.CPU), proto.CommodityDTO_VCPU, resizeable)
	result = append(result, cpuComm)

	memComm, _ := CreateResourceCommodityResize(&(docker.Memory), proto.CommodityDTO_VMEM, resizeable)
	result = append(result, memComm)

	appComm, _ := CreateKeyCommodity(docker.UUID, proto.CommodityDTO_APPLICATION)
	result = append(result, appComm)

	return result, nil
}
