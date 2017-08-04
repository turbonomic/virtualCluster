package target

import (
	"fmt"
	"github.com/golang/glog"
)

type ClusterBuilder struct {
	clusterId   string
	clusterName string

	topology *TargetTopology

	containers map[string]*Container
	pods       map[string]*Pod
	nodes      []*HostNode
	services   []*VirtualApp
}

func NewClusterBuilder(clusterId, clusterName, topoConf string) *ClusterBuilder {
	topo := NewTargetTopology(clusterId)
	if err := topo.LoadTopology(topoConf); err != nil {
		glog.Errorf("failed to load topology from file: %s, error: %v",
			topoConf, err)
		return nil
	}

	return &ClusterBuilder{
		clusterId:   clusterId,
		clusterName: clusterName,
		topology:    topo,

		//containers: make(map[string]*Container),
		//pods: make(map[string]*Pod),
		//nodes: []*HostNode{},
		//services: []*VirtualApp{},
	}
}

func (b *ClusterBuilder) buildContainers() error {
	containers := make(map[string]*Container)

	for k, v := range b.topology.ContainerTemplateMap {
		container := NewContainer(k, k)

		container.CPU = v.CPU
		container.Memory = v.Memory

		containers[k] = container
	}

	b.containers = containers
	return nil
}

func (b *ClusterBuilder) buildPods() error {
	result := make(map[string]*Pod)

	allContainers := b.containers

	for k, v := range b.topology.PodTemplateMap {
		pod := NewPod(k, k)

		containers := []*Container{}
		for i, cname := range v.Containers {
			if container, exist := allContainers[cname]; exist {
				// generate a new container with different UUID
				newId := fmt.Sprintf("%s-%d", pod.UUID, i)
				ct := container.Clone(newId)
				containers = append(containers, ct)
			} else {
				glog.Warningf("pod[%s]-%dth container[%s] does not exist.", k, i+1, cname)
				break
			}
		}

		glog.V(3).Infof("pod[%s] has %d containers.", k, len(containers))

		if len(containers) < 1 {
			glog.Warningf("pod[%s] has no container.", k)
			continue
		}

		if len(containers) != len(v.Containers) {
			glog.Warningf("cannot get enough containers[%d Vs. %d] for pod[%s].",
				len(containers), len(v.Containers), k)
			continue
		}

		pod.Containers = containers
		result[k] = pod
	}

	b.pods = result
	return nil
}

func (b *ClusterBuilder) buildNodes() error {
	var result []*HostNode

	allPods := b.pods
	for k, v := range b.topology.NodeTemplateMap {
		node := NewHostNode(k, k)
		node.ClusterID = b.clusterId
		//node.IP = ip

		pods := []*Pod{}
		for i, podName := range v.Pods {
			if pod, exist := allPods[podName]; exist {
				pods = append(pods, pod)
			} else {
				glog.Warningf("node[%s]-%dth pod[%s] does not exist.", k, i+1, podName)
				break
			}
		}

		glog.V(3).Infof("node[%s] has %d Pods.", k, len(pods))
		if len(pods) != len(v.Pods) {
			glog.Warningf("cannot get enough Pods[%d Vs. %d] for node[%s].",
				len(pods), len(v.Pods), k)
			continue
		}

		node.Pods = pods
		result = append(result, node)
	}

	b.nodes = result
	return nil
}

func (b *ClusterBuilder) buildVirtualApp() error {
	var result []*VirtualApp

	allPods := b.pods
	for k, v := range b.topology.ServiceTemplateMap {
		vapp := NewVirtualApp(k, k)

		pods := []*Pod{}
		for i, podName := range v.Pods {
			if pod, exist := allPods[podName]; exist {
				pods = append(pods, pod)
			} else {
				glog.Warningf("vapp[%s]-%dth pod[%s] does not exist.", k, i+1, podName)
				break
			}
		}

		glog.V(3).Infof("vapp[%s] has %d Pods.", k, len(pods))
		if len(pods) != len(v.Pods) {
			glog.Warningf("cannot get enough Pods[%d Vs. %d] for vapp[%s].",
				len(pods), len(v.Pods), k)
			continue
		}

		vapp.Pods = pods
		result = append(result, vapp)
	}

	b.services = result
	return nil
}

func (b *ClusterBuilder) GenerateCluster() (*Cluster, error) {
	if b.topology == nil {
		err := fmt.Errorf("need to set topology first.")
		glog.Error(err.Error())
		return nil, err
	}

	if err := b.buildContainers(); err != nil {
		err := fmt.Errorf("Generate cluster failed: build containers failed: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	if err := b.buildPods(); err != nil {
		err := fmt.Errorf("Generate cluster failed: build pods failed: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	if err := b.buildNodes(); err != nil {
		err := fmt.Errorf("Generate cluster failed: build pods failed: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	if err := b.buildVirtualApp(); err != nil {
		err := fmt.Errorf("Generate cluster failed: build pods failed: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	cluster := NewCluster(b.clusterName, b.clusterId)
	cluster.Nodes = b.nodes
	cluster.Services = b.services

	cluster.GenerateContainerAPP()
	return cluster, nil
}
