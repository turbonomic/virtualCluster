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
	vnodes     map[string]*VNode
	nodes      map[string]*Node
	services   []*VirtualApp
}

func NewClusterBuilderfromTopology(clusterId, clusterName string, topo *TargetTopology) *ClusterHandler {
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
	containers := make(map[string]*Container)

	for k, v := range b.topology.ContainerTemplateMap {
		container := NewContainer(k, k)

		container.CPU = v.CPU
		container.Memory = v.Memory

		containers[k] = container
		glog.V(4).Infof("container-%+v", container)
	}

	b.containers = containers
	return nil
}

func resetResource(r *Resource) {
	r.Capacity = 0
	r.Used = 0
}

func addResource(r, delta *Resource) {
	r.Capacity += delta.Capacity
	r.Used += delta.Used
}

func setResource(r, r2 *Resource) {
	r.Capacity = r2.Capacity
	r.Used = r2.Used
}

func (b *ClusterBuilder) buildPods() error {
	result := make(map[string]*Pod)

	allContainers := b.containers

	cpu := &Resource{}
	mem := &Resource{}

	for k, v := range b.topology.PodTemplateMap {
		pod := NewPod(k, k)

		resetResource(cpu)
		resetResource(mem)

		containers := []*Container{}
		for i, cname := range v.Containers {
			if container, exist := allContainers[cname]; exist {
				// generate a new container with different UUID
				newId := fmt.Sprintf("%s-%s", container.Name, pod.UUID)
				ct := container.Clone(newId, newId)
				containers = append(containers, ct)
				addResource(cpu, &(ct.CPU))
				addResource(mem, &(ct.Memory))
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
		setResource(&(pod.CPU), cpu)
		setResource(&(pod.Memory), mem)
		result[k] = pod
		glog.V(4).Infof("pod--%+v", pod)
	}

	b.pods = result
	return nil
}

func assignVNode(node *VNode, tmp *vnodeTemplate) {
	node.Memory.Capacity = tmp.Memory
	node.CPU.Capacity = tmp.CPU
	node.IP = tmp.IP
}

func (b *ClusterBuilder) buildVNodes() error {
	result := make(map[string]*VNode)

	allPods := b.pods
	for k, v := range b.topology.VNodeTemplateMap {
		vnode := NewVNode(k, k)
		assignVNode(vnode, v)
		vnode.ClusterId = b.clusterId

		cpu := 0.0
		mem := 0.0

		pods := make(map[string]*Pod)
		for i, podName := range v.Pods {
			if pod, exist := allPods[podName]; exist {
				pods[pod.UUID] = pod
				cpu += pod.CPU.Used
				mem += pod.Memory.Used
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
		vnode.CPU.Used = cpu
		vnode.Memory.Used = mem

		result[vnode.UUID] = vnode
		glog.V(4).Infof("[vnode] %+v", vnode)
	}

	b.vnodes = result
	return nil
}

func assignNode(node *Node, tmp *nodeTemplate) {
	node.Memory.Capacity = tmp.Memory
	node.CPU.Capacity = tmp.CPU
	node.IP = tmp.IP
}

func (b *ClusterBuilder) buildNodes() error {
	result := make(map[string]*Node)

	allVMs := b.vnodes
	for k, v := range b.topology.NodeTemplateMap {
		node := NewNode(k, k)
		assignNode(node, v)
		node.ClusterId = b.clusterId

		cpu := 0.0
		mem := 0.0

		vnodes := make(map[string]*VNode)
		for i, vmKey := range v.VMs {
			if vm, exist := allVMs[vmKey]; exist {
				vnodes[vm.UUID] = vm
				cpu += vm.CPU.Capacity
				mem += vm.Memory.Capacity
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
		node.CPU.Used = cpu
		node.Memory.Used = mem

		result[node.UUID] = node
		glog.V(4).Infof("[node] %+v", node)
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

	cluster := NewCluster(b.clusterName, b.clusterId)
	cluster.Nodes = b.nodes
	cluster.Services = b.services

	cluster.CompleteBuild()
	return cluster, nil
}
