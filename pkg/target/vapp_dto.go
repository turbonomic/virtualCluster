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

	err, qps, rt := vapp.getCommoditiesBought(vAppBuilder)
	if err != nil {
		nerr := fmt.Errorf("build VApplication DTO failed: %v", err)
		glog.Error(nerr.Error())
		return nil, nerr
	}

	var commsSold []*proto.CommodityDTO
	transComm, _ := CreateCapacityUsedCommodity(vapp.UUID, &Resource{Used: qps.Used, Capacity: qps.Capacity},
		proto.CommodityDTO_TRANSACTION)
	rtComm, _ := CreateCapacityUsedCommodity(vapp.UUID, &Resource{Used: rt.Used, Capacity: rt.Capacity},
		proto.CommodityDTO_RESPONSE_TIME)
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
	used float64, capacity float64) (*proto.CommodityDTO, error) {
	app := container.App
	appCommodity, err := builder.NewCommodityDTOBuilder(commodityType).
		Key(app.UUID).
		Used(used).
		Capacity(capacity).
		Create()
	if err != nil {
		glog.Errorf("failed to create commodity bought for VirtualApp[%s]-pod[%s]-container[%s]-app[%s]",
			vapp.Name, pod.Name, container.Name, app.Name)
	}
	return appCommodity, err
}

func (vapp *VirtualApp) getCommoditiesBought(vAppBuilder *builder.EntityDTOBuilder) (error, Resource, Resource) {
	i := 0
	qpsused := 0.0
	qpscap := 0.0
	rtused := 0.0
	rtcap := 0.0
	weightedRT := 0.0

	for _, pod := range vapp.Pods {
		for _, container := range pod.Containers {
			if container.App == nil {
				glog.Errorf("contain.App is not ready; VirtualApp[%s]-pod[%s]-container[%s]",
					vapp.Name, pod.Name, container.Name)
				continue
			}

			app := container.App
			appCommodity, err := vapp.createAppCommodity(pod, container, proto.CommodityDTO_TRANSACTION,
				app.QPS.Used, app.QPS.Capacity)
			if err != nil {
				continue
			}
			qpsused += app.QPS.Used
			qpscap += app.QPS.Capacity
			weightedRT += (app.QPS.Used * app.ResponseTime.Used)
			bought := []*proto.CommodityDTO{appCommodity}

			appCommodity, err = vapp.createAppCommodity(pod, container, proto.CommodityDTO_RESPONSE_TIME,
				app.ResponseTime.Used, app.ResponseTime.Capacity)
			if err == nil {
				rtused += app.ResponseTime.Used
				rtcap += app.ResponseTime.Capacity
				bought = append(bought, appCommodity)
			}

			appProvider := builder.CreateProvider(proto.EntityDTO_APPLICATION, app.UUID)
			vAppBuilder.Provider(appProvider).BuysCommodities(bought)
			i++
		}
	}

	if i < 1 {
		return fmt.Errorf("cannot get commodities bought from containers for VApp[%s/%s]", vapp.Kind, vapp.Name),
			Resource{Capacity: 0.0, Used: 0.0}, Resource{Capacity: 0.0, Used: 0.0}
	}

	fi := float64(i)
	var weightedResponseTime float64
	if qpsused == 0.0 {
		// We cannot calculate a weighted response time, so use the average instead
		weightedResponseTime = rtused / fi
	} else {
		weightedResponseTime = weightedRT / qpsused
	}
	return nil, Resource{Capacity: qpscap, Used: qpsused},
		Resource{Capacity: rtcap / fi, Used: weightedResponseTime}
}
