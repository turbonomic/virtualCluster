package target

import (
	"fmt"
	"github.com/golang/glog"
	"sync"

	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type ClusterHandler struct {
	cluster *Cluster

	containers map[string]*Container
	pods       map[string]*Pod
	vnodes     map[string]*VNode
	nodes      map[string]*Node
	switches   map[string]*Switch

	Ready bool
	mux   sync.Mutex
}

func NewClusterHandler(c *Cluster) *ClusterHandler {

	h := &ClusterHandler{
		cluster: c,
		Ready:   false,
	}

	h.BuildIndex()

	return h
}

func (h *ClusterHandler) String() string {
	return fmt.Sprintf("clusterInfo: %s; %s", h.cluster.Name, h.cluster.UUID)
}

func (h *ClusterHandler) BuildIndex() {
	h.mux.Lock()
	defer h.mux.Unlock()
	if h.Ready {
		glog.V(3).Infof("Index is ready.")
		return
	}

	containers := make(map[string]*Container)
	pods := make(map[string]*Pod)
	vnodes := make(map[string]*VNode)
	nodes := make(map[string]*Node)
	switches := make(map[string]*Switch)

	c := h.cluster
	for _, host := range c.Nodes {
		nodes[host.UUID] = host

		for _, vhost := range host.VMs {
			vnodes[vhost.UUID] = vhost

			for _, pod := range vhost.Pods {
				pods[pod.UUID] = pod

				for _, container := range pod.Containers {
					containers[container.UUID] = container
				}
			}
		}
	}

	h.switches = switches
	h.nodes = nodes
	h.vnodes = vnodes
	h.pods = pods
	h.containers = containers

	h.Ready = true
}

func (h *ClusterHandler) GenerateClusterDTOs() ([]*proto.EntityDTO, error) {
	h.mux.Lock()
	defer h.mux.Unlock()
	return h.cluster.GenerateDTOs()
}

func (h *ClusterHandler) MovePod(podId, vnodeId string) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	if !h.Ready {
		err := fmt.Errorf("ClusterHandler is not ready.")
		glog.Error(err.Error())
		return err
	}

	//0. check
	pod, exist := h.pods[podId]
	if !exist {
		err := fmt.Errorf("MovePod failed. Pod[%s] is not found", podId)
		glog.Error(err.Error())
		return err
	}

	if pod.ProviderID == vnodeId {
		msg := fmt.Sprintf("MovePod aborted. Pod[%s][%s] is already on vnode[%s].", podId, pod.Name, vnodeId)
		glog.Warning(msg)
		return nil
	}

	vnode, exist := h.vnodes[vnodeId]
	if !exist {
		err := fmt.Errorf("MovePod failed. VNode[%s] is not found", vnodeId)
		glog.Error(err.Error())
		return err
	}

	//1. delete it from original VNode
	oldVnode, exist := h.vnodes[pod.ProviderID]
	if !exist {
		err := fmt.Errorf("MovePod failed. Cannot found original VNode[%s].", pod.ProviderID)
		glog.Error(err.Error())
		return err
	}
	if err := oldVnode.DeletePod(podId); err != nil {
		err := fmt.Errorf("MovePod failed. %v", err)
		glog.Error(err.Error())
		return err
	}

	//2. add it to the new VNode
	if err := vnode.AddPod(pod); err != nil {
		err := fmt.Errorf("MovePod failed. %v", err)
		glog.Error(err.Error())
		return err
	}

	glog.V(2).Infof("Successed: move pod[%s] from vnode[%s] to vnode[%s]", pod.Name, oldVnode.Name, vnode.Name)
	glog.V(2).Infof("oldVnode pods: %s", oldVnode.GetPodNames())
	glog.V(2).Infof("newVnode pods: %s", vnode.GetPodNames())

	return nil
}

func (h *ClusterHandler) ResizeContainerCapacity(containerId string, cpu, memory float64) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	if !h.Ready {
		err := fmt.Errorf("ClusterHandler is not ready.")
		glog.Error(err.Error())
		return err
	}

	container, exist := h.containers[containerId]
	if !exist {
		err := fmt.Errorf("ResizeContainerCapacity failed. container[%s] is not found.", containerId)
		glog.Error(err.Error())
		return err
	}

	container.SetCapacity(cpu, memory)

	return nil
}

func (h *ClusterHandler) MoveVirtualMachine(vnodeId, nodeId string) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	if !h.Ready {
		err := fmt.Errorf("ClusterHandler is not ready.")
		glog.Error(err.Error())
		return err
	}

	//0. check
	vnode, exist := h.vnodes[vnodeId]
	if !exist {
		err := fmt.Errorf("MoveVM failed. VirtualMachine[%s] is not found", vnodeId)
		glog.Error(err.Error())
		return err
	}

	if vnode.ProviderID == nodeId {
		msg := fmt.Sprintf("MoveVM aborted. VM[%s][%s] is already on node[%s].", vnode.Name, vnodeId, nodeId)
		glog.Warning(msg)
		return nil
	}

	node, exist := h.nodes[nodeId]
	if !exist {
		err := fmt.Errorf("MoveVM failed. New destion Node[%s] is not found", nodeId)
		glog.Error(err.Error())
		return err
	}

	//1. delete it from original VNode
	oldNode, exist := h.nodes[vnode.ProviderID]
	if !exist {
		err := fmt.Errorf("MoveVM failed. Cannot found original Node[%s].", vnode.ProviderID)
		glog.Error(err.Error())
		return err
	}
	if err := oldNode.DeleteVM(vnodeId); err != nil {
		err := fmt.Errorf("MovePod failed. %v", err)
		glog.Error(err.Error())
		return err
	}

	//2. add it to the new VNode
	if err := node.AddVM(vnode); err != nil {
		err := fmt.Errorf("MovePod failed. %v", err)
		glog.Error(err.Error())
		return err
	}

	glog.V(2).Infof("Successed: move vnode[%s] from node[%s] to node[%s]", vnode.Name, oldNode.Name, node.Name)
	glog.V(2).Infof("old Node has vnodes: %s", oldNode.GetVMNames())
	glog.V(2).Infof("new Node has vnodes: %s", node.GetVMNames())

	return nil
}
