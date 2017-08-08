package target

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func (c *Cluster) GenerateDTOs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	if c.Nodes == nil || len(c.Nodes) < 1 {
		err := fmt.Errorf("empty cluster[%s/%s].", c.Name, c.UUID)
		glog.Error(err.Error())
		return result, err
	}

	//1. node, pod, container, app DTOs
	for _, host := range c.Nodes {
		hostDTO, err := host.BuildDTO()
		if err != nil {
			e := fmt.Errorf("failed to build hostDTO for node[%s]", host.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, hostDTO)

		subDTOs, err := host.BuildSubDTOs()
		if err != nil {
			e := fmt.Errorf("failed to build subHostDTOs for node[%s]", host.Name)
			glog.Error(e.Error())
			continue
		}
		result = append(result, subDTOs...)
		glog.V(3).Infof("There are %d DTOs on node[%s].", len(subDTOs)+1, host.Name)
	}

	//2. service DTOs
	if serviceDTOs, err := c.generateServiceDTOs(); err != nil {
		glog.Errorf("failed to generate ServiceDTOs:%v", err)
	} else {
		result = append(result, serviceDTOs...)
	}

	glog.V(2).Infof("There are %d DTOs in total.", len(result))
	if len(result) < 1 {
		return result, fmt.Errorf("failed to generate valid DTOs.")
	}

	return result, nil
}

func (c *Cluster) generateServiceDTOs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO
	if c.Services == nil || len(c.Services) < 1 {
		glog.Warningf("No services in cluster[%s]", c.Name)
		return result, nil
	}

	for _, service := range c.Services {
		serviceDTO, err := service.BuildDTO()
		if err != nil {
			e := fmt.Errorf("failed to build serviceDTO for service[%s]: %v", service.Name, err)
			glog.Error(e.Error())
			continue
		}

		result = append(result, serviceDTO)
	}

	glog.V(2).Infof("There are %d services, and %d serviceDTOs.", len(c.Services), len(result))
	return result, nil
}

// Generate complement information
//   (1) set providerId for each entity;
//   (2) Generate Application Entity for each container;
func (c *Cluster) CompleteBuild() {
	for _, host := range c.Nodes {
		for _, vhost := range host.VMs {
			vhost.ProviderID = host.UUID
			for _, pod := range vhost.Pods {
				pod.ProviderID = vhost.UUID
				for _, container := range pod.Containers {
					container.ProviderID = pod.UUID
					container.GenerateApp()
				}
			}
		}
	}
}
