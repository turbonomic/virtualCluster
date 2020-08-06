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

	//0. calculate the resource usage
	c.SetResourceAmount()

	//1. switch, node, pod, container, app DTOs
	if c.Switches != nil {
		for _, networkswitch := range c.Switches {
			switchDTO, err := networkswitch.BuildDTO()
			if err != nil {
				e := fmt.Errorf("failed to build switchDTO for switch[%s]", networkswitch.Name)
				glog.Error(e.Error())
				continue
			}
			result = append(result, switchDTO)

			subDTOs, err := networkswitch.BuildSubDTOs()
			if err != nil {
				e := fmt.Errorf("failed to build subHostDTOs for node[%s]", networkswitch.Name)
				glog.Error(e.Error())
				continue
			}
			result = append(result, subDTOs...)
			glog.V(3).Infof("There are %d DTOs on node[%s].", len(subDTOs)+1, networkswitch.Name)
		}
	} else {
		for _, host := range c.Nodes {
			hostDTO, err := host.BuildDTO(nil)
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

	glog.V(3).Infof("There are %d services, and %d serviceDTOs.", len(c.Services), len(result))
	return result, nil
}

// 1. set ProviderId for each SE;
// 2. Generate Application for each pod-container;
// 3. calculate and set resource usage;
func (c *Cluster) CompleteBuild() {
	c.SetProvider()
	c.SetResourceAmount()
}

// Generate complement information
//   (1) set providerId for each entity;
//   (2) Generate Application Entity for each container;
func (c *Cluster) SetProvider() {
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
	return
}

// SetResourceAmount: Set the resource Capacity and Usage
// Container.Capacity = Container.Limit/Pod.Capacity
// Pod.Capacity = VM.Capacity
// VM.Capacity = setting
// PM.Capacity = setting
// Container.Used = monitored (from topology)
// Pod.Used = sum.container.Used
// VM.Used = monitored = sum.Pod.Used + overhead1
// PM.Used = monitored = sum.Vm.Used + overhead2
func (c *Cluster) SetResourceAmount() {
	for _, host := range c.Nodes {
		hostCPU := 0.0
		hostMem := 0.0

		for _, vhost := range host.VMs {
			vhostCPU := 0.0
			vhostMem := 0.0

			for _, pod := range vhost.Pods {
				pod.CPU.Capacity = vhost.CPU.Capacity
				pod.Memory.Capacity = vhost.Memory.Capacity

				podCPU := 0.0
				podMem := 0.0

				for _, container := range pod.Containers {
					app := container.App
					app.CPU.Used = container.CPU.Used
					app.Memory.Used = container.Memory.Used

					podCPU += container.CPU.Used
					podMem += container.Memory.Used

					if container.CPU.Capacity < 1 {
						container.CPU.Capacity = pod.CPU.Capacity
					}

					if container.Memory.Capacity < 1 {
						container.Memory.Capacity = pod.Memory.Capacity
					}
				}

				pod.CPU.Used = podCPU
				pod.Memory.Used = podMem

				vhostCPU += pod.CPU.Used
				vhostMem += pod.Memory.Used
			}

			vhost.CPU.Used = vhostCPU + defaultOverheadVMCPU
			vhost.Memory.Used = vhostMem + defaultOverheadVMMem

			hostCPU += vhost.CPU.Used
			hostMem += vhost.Memory.Used
		}

		host.CPU.Used = hostCPU + defaultOverheadPMCPU
		host.Memory.Used = hostMem + defaultOverheadPMMem
	}

	return
}
