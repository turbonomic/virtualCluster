package topology

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/virtualCluster/pkg/target"
)

type ClusterBuilder struct {
	clusterId   string
	clusterName string

	topology *TargetTopology

	containers map[string]*target.Container
	pods       map[string]*target.Pod
	vnodes     map[string]*target.VNode
	nodes      map[string]*target.Node
	services   []*target.VirtualApp
}

func NewClusterBuilderfromTopology(clusterId, clusterName string, topo *TargetTopology) *ClusterBuilder {
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

func NewClusterBuilder(clusterId, clusterName, topoConf string) *ClusterBuilder {
	topo := NewTargetTopology(clusterId)
	if err := topo.LoadTopology(topoConf); err != nil {
		glog.Errorf("failed to load topology from file: %s, error: %v",
			topoConf, err)
		return nil
	}

	return NewClusterBuilderfromTopology(clusterId, clusterName, topo)
}

func (b *ClusterBuilder) buildContainers() error {
	containers := make(map[string]*target.Container)

	for k, v := range b.topology.ContainerTemplateMap {
		container := target.NewContainer(k, k)

		container.CPU = v.CPU
		container.Memory = v.Memory
		container.ReqCPU = v.ReqCPU
		container.ReqMemory = v.ReqMem
		container.QPS = v.QPS
		container.ResponseTime = v.ResponseTime

		containers[k] = container
		glog.V(4).Infof("container-%+v", container)
	}

	b.containers = containers
	return nil
}

//Note: will set Pod resource amount later in cluster.SetResourceAmount()
func (b *ClusterBuilder) buildPods() error {
	result := make(map[string]*target.Pod)

	allContainers := b.containers

	for k, v := range b.topology.PodTemplateMap {
		pod := target.NewPod(k, k)

		containers := []*target.Container{}
		for i, cname := range v.Containers {
			if container, exist := allContainers[cname]; exist {
				// generate a new container with different UUID
				newId := fmt.Sprintf("%s-%s", container.Name, pod.UUID)
				ct := container.Clone(newId, newId)
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
		glog.V(4).Infof("pod--%+v", pod)
	}

	b.pods = result
	return nil
}

func assignVNode(node *target.VNode, tmp *vnodeTemplate) {
	node.Memory.Capacity = tmp.Memory
	node.CPU.Capacity = tmp.CPU
	node.IP = tmp.IP
}

//Note: will set VNode resourceAmount in cluster.SetResourceAmount()
func (b *ClusterBuilder) buildVNodes() error {
	result := make(map[string]*target.VNode)

	allPods := b.pods
	for k, v := range b.topology.VNodeTemplateMap {
		vnode := target.NewVNode(k, k)
		assignVNode(vnode, v)
		vnode.ClusterId = b.clusterId

		pods := make(map[string]*target.Pod)
		for i, podName := range v.Pods {
			if pod, exist := allPods[podName]; exist {
				pods[pod.UUID] = pod
			} else {
				glog.Warningf("vnode[%s]-%dth pod[%s] does not exist.", k, i+1, podName)
				break
			}
		}

		glog.V(3).Infof("vnode[%s] has %d Pods.", k, len(pods))
		if len(pods) != len(v.Pods) {
			glog.Warningf("cannot get enough Pods[%d Vs. %d] for vnode[%s].",
				len(pods), len(v.Pods), k)
			continue
		}

		vnode.Pods = pods
		result[vnode.UUID] = vnode
		glog.V(4).Infof("[vnode] %+v", vnode)
	}

	b.vnodes = result
	return nil
}

func assignNode(node *target.Node, tmp *nodeTemplate) {
	node.Memory.Capacity = tmp.Memory
	node.CPU.Capacity = tmp.CPU
	node.IP = tmp.IP
}

//Note: will set Node resourceAmount in cluster.SetResourceAmount()
func (b *ClusterBuilder) buildNodes() error {
	result := make(map[string]*target.Node)

	allVMs := b.vnodes
	for k, v := range b.topology.NodeTemplateMap {
		node := target.NewNode(k, k)
		assignNode(node, v)
		node.ClusterId = b.clusterId

		vnodes := make(map[string]*target.VNode)
		for i, vmKey := range v.VMs {
			if vm, exist := allVMs[vmKey]; exist {
				vnodes[vm.UUID] = vm
			} else {
				glog.Warningf("node[%s]-%dth VM[%s] does not exist.", k, i+1, vmKey)
				break
			}
		}

		glog.V(3).Infof("node[%s] has %d VNodes.", k, len(vnodes))
		if len(vnodes) != len(v.VMs) {
			glog.Warningf("cannot get enough VMs[%d Vs. %d] for node[%s].",
				len(vnodes), len(v.VMs), k)
			continue
		}

		node.VMs = vnodes
		result[node.UUID] = node
		glog.V(4).Infof("[node] %+v", node)
	}

	b.nodes = result
	return nil
}

func (b *ClusterBuilder) buildVirtualApp() error {
	var result []*target.VirtualApp

	allPods := b.pods
	for k, v := range b.topology.ServiceTemplateMap {
		vapp := target.NewVirtualApp(k, k)

		pods := []*target.Pod{}
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

func (b *ClusterBuilder) GenerateCluster() (*target.Cluster, error) {
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

	if err := b.buildVNodes(); err != nil {
		err := fmt.Errorf("Generate cluster failed: build vnodes failed: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	if err := b.buildNodes(); err != nil {
		err := fmt.Errorf("Generate cluster failed: build nodes failed: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	if err := b.buildVirtualApp(); err != nil {
		err := fmt.Errorf("Generate cluster failed: build virtualApp failed: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	cluster := target.NewCluster(b.clusterName, b.clusterId)
	cluster.Nodes = b.nodes
	cluster.Services = b.services

	cluster.CompleteBuild()
	return cluster, nil
}
