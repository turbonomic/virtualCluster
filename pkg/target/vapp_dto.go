package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (vapp *VirtualApp) BuildDTO() (*proto.EntityDTO, error) {
	vAppBuilder := builder.
		NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, vapp.UUID).
		DisplayName(vapp.Name)

	if err := vapp.getCommoditiesBought(vAppBuilder); err != nil {
		nerr := fmt.Errorf("build VApplication DTO failed: %v", err)
		glog.Error(nerr.Error())
		return nil, nerr
	}

	// For now, just add dummy Transaction and ResponseTime commodities.
	// TODO Only add a commodity if at least one underlying Application provides it.
	dummyResource := &Resource{Used: 0.0, Capacity: 0.0}
	var commsSold []*proto.CommodityDTO
	transComm, _ := CreateTransactionCommodity(vapp.UUID, dummyResource, proto.CommodityDTO_TRANSACTION)
	rtComm, _ := CreateResponseTimeCommodity(vapp.UUID, dummyResource, proto.CommodityDTO_RESPONSE_TIME)
	commsSold = append(commsSold, rtComm, transComm)
	vAppBuilder.SellsCommodities(commsSold)

	vappData := &proto.EntityDTO_VirtualApplicationData{
		ServiceType: &(vapp.Name),
	}
	vAppBuilder.VirtualApplicationData(vappData)
	vAppBuilder.WithPowerState(proto.EntityDTO_POWERED_ON)

	entity, err := vAppBuilder.Create()
	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for vApplication(%v): %v",
			vapp.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (vapp *VirtualApp) createAppCommodity(pod *Pod, container *Container,
		commodityType proto.CommodityDTO_CommodityType,
		used float64) (*proto.CommodityDTO, error) {
	app := container.App
	appCommodity, err := builder.NewCommodityDTOBuilder(commodityType).
		Key(app.UUID).
		Used(used).
		Create()
	if err != nil {
		glog.Errorf("failed to create commodity bought for VirtualApp[%s]-pod[%s]-container[%s]-app[%s]",
			vapp.Name, pod.Name, container.Name, app.Name)
	}
	return appCommodity, err
}

func (vapp *VirtualApp) getCommoditiesBought(vAppBuilder *builder.EntityDTOBuilder) error {
	i := 0

	for _, pod := range vapp.Pods {
		for _, container := range pod.Containers {
			if container.App == nil {
				glog.Errorf("contain.App is not ready; VirtualApp[%s]-pod[%s]-container[%s]",
					vapp.Name, pod.Name, container.Name)
				continue
			}

			app := container.App
			appCommodity, err := vapp.createAppCommodity(pod, container, proto.CommodityDTO_TRANSACTION, app.QPS.Used)
			if err != nil {
				continue
			}
			bought := []*proto.CommodityDTO{appCommodity}

			if app.ResponseTime != nil {
				appCommodity, err := vapp.createAppCommodity(pod, container, proto.CommodityDTO_RESPONSE_TIME,
					app.ResponseTime.Used)
				if err == nil {
					bought = append(bought, appCommodity)
				}
			}

			appProvider := builder.CreateProvider(proto.EntityDTO_APPLICATION, app.UUID)
			vAppBuilder.Provider(appProvider).BuysCommodities(bought)
			i++
		}
	}

	if i < 1 {
		return fmt.Errorf("cannot get commodities bought from containers for VApp[%s/%s]", vapp.Kind, vapp.Name)
	}

	return nil
}
