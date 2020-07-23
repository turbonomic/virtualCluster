package target

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (service *VirtualApp) BuildDTO() (*proto.EntityDTO, error) {
	serviceBuilder := builder.
		NewEntityDTOBuilder(proto.EntityDTO_SERVICE, service.UUID).
		DisplayName(service.Name)

	err, qps, rt := service.getCommoditiesBought(serviceBuilder)
	if err != nil {
		nerr := fmt.Errorf("build VApplication DTO failed: %v", err)
		glog.Error(nerr.Error())
		return nil, nerr
	}

	var commsSold []*proto.CommodityDTO
	transComm, _ := CreateCapacityUsedCommodity(service.UUID, &Resource{Used: qps.Used, Capacity: qps.Capacity},
		proto.CommodityDTO_TRANSACTION)
	rtComm, _ := CreateCapacityUsedCommodity(service.UUID, &Resource{Used: rt.Used, Capacity: rt.Capacity},
		proto.CommodityDTO_RESPONSE_TIME)
	commsSold = append(commsSold, rtComm, transComm)
	serviceBuilder.SellsCommodities(commsSold)

	serviceData := &proto.EntityDTO_ServiceData{
		ServiceType: &service.Name,
	}
	serviceBuilder.ServiceData(serviceData)
	serviceBuilder.WithPowerState(proto.EntityDTO_POWERED_ON)

	entity, err := serviceBuilder.Create()
	if err != nil {
		msg := fmt.Errorf("Failed to build EntityDTO for vApplication(%v): %v",
			service.Name, err.Error())
		glog.Error(msg.Error())
		return nil, msg
	}

	return entity, nil
}

func (service *VirtualApp) createAppCommodity(pod *Pod, container *Container,
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
			service.Name, pod.Name, container.Name, app.Name)
	}
	return appCommodity, err
}

func (service *VirtualApp) getCommoditiesBought(vAppBuilder *builder.EntityDTOBuilder) (error, Resource, Resource) {
	i := 0
	qpsused := 0.0
	qpscap := 0.0
	rtused := 0.0
	rtcap := 0.0
	weightedRT := 0.0

	for _, pod := range service.Pods {
		for _, container := range pod.Containers {
			if container.App == nil {
				glog.Errorf("contain.App is not ready; VirtualApp[%s]-pod[%s]-container[%s]",
					service.Name, pod.Name, container.Name)
				continue
			}

			app := container.App
			appCommodity, err := service.createAppCommodity(pod, container, proto.CommodityDTO_TRANSACTION,
				app.QPS.Used, app.QPS.Capacity)
			if err != nil {
				continue
			}
			qpsused += app.QPS.Used
			qpscap += app.QPS.Capacity
			weightedRT += (app.QPS.Used * app.ResponseTime.Used)
			bought := []*proto.CommodityDTO{appCommodity}

			appCommodity, err = service.createAppCommodity(pod, container, proto.CommodityDTO_RESPONSE_TIME,
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
		return fmt.Errorf("cannot get commodities bought from containers for VApp[%s/%s]", service.Kind, service.Name),
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
