package target

import (
	"fmt"
	"github.com/golang/glog"
)

const (
	KindApp        = "application"
	KindContainer  = "container"
	KindPod        = "pod"
	KindVirtualApp = "service"
	KindVNode      = "vhost"
	KindNode       = "host"
	KindCluster    = "cluster"

	emptyProvider = "None"
)

type ObjectMeta struct {
	Name       string
	UUID       string
	Kind       string
	ProviderID string
}

type Resource struct {
	Capacity float64
	Used     float64
}

type Application struct {
	ObjectMeta

	CPU         Resource
	Memory      Resource
	Transaction float64
}

type Container struct {
	ObjectMeta

	CPU       Resource
	Memory    Resource
	ReqCPU    float64
	ReqMemory float64

	App *Application
}

type Pod struct {
	ObjectMeta
	//AppName string

	CPU    Resource
	Memory Resource

	Containers []*Container
}

type VirtualApp struct {
	ObjectMeta

	Pods []*Pod
}

// virtual machine
type VNode struct {
	ObjectMeta

	CPU    Resource
	Memory Resource

	ClusterId string
	IP        string

	//a map for easy of move/deletion, key=pod.UUID
	Pods map[string]*Pod
}

// physical machine
type Node struct {
	ObjectMeta

	CPU    Resource
	Memory Resource

	ClusterId string
	IP        string

	//Map for easy of deletion
	// key = vm.UUID
	VMs map[string]*VNode
}

type Cluster struct {
	ObjectMeta
	Nodes    map[string]*Node
	Services []*VirtualApp
}

func NewContainer(name, id string) *Container {
	return &Container{
		ObjectMeta: ObjectMeta{
			Kind: KindContainer,
			Name: name,
			UUID: id,
		},
	}
}

func NewPod(name, id string) *Pod {
	return &Pod{
		ObjectMeta: ObjectMeta{
			Kind: KindPod,
			Name: name,
			UUID: id,
		},
	}
}

func NewVNode(name, id string) *VNode {
	return &VNode{
		ObjectMeta: ObjectMeta{
			Kind: KindVNode,
			Name: name,
			UUID: id,
		},
	}
}

func NewNode(name, id string) *Node {
	return &Node{
		ObjectMeta: ObjectMeta{
			Kind: KindNode,
			Name: name,
			UUID: id,
		},
	}
}

func NewApplication(name, id string) *Application {
	return &Application{
		ObjectMeta: ObjectMeta{
			Kind: KindApp,
			Name: name,
			UUID: id,
		},
	}
}

func NewVirtualApp(name, id string) *VirtualApp {
	return &VirtualApp{
		ObjectMeta: ObjectMeta{
			Kind: KindVirtualApp,
			Name: name,
			UUID: id,
		},
	}
}

func NewCluster(name, id string) *Cluster {
	return &Cluster{
		ObjectMeta: ObjectMeta{
			Kind: KindCluster,
			Name: name,
			UUID: id,
		},
	}
}

func (v *VNode) DeletePod(podId string) error {
	pod, exist := v.Pods[podId]

	if !exist {
		err := fmt.Errorf("VNode Delete Pod failed: VNode[%s][%s] does not has pod[%s].", v.Name, v.UUID, podId)
		glog.Error(err.Error())
		return err
	}

	pod.ProviderID = emptyProvider
	delete(v.Pods, podId)

	return nil
}

func (v *VNode) AddPod(pod *Pod) error {
	podId := pod.UUID

	if _, exist := v.Pods[podId]; exist {
		err := fmt.Errorf("VNode Add Pod failed: VNode[%s][%s] already has pod[%s].", v.Name, v.UUID, podId)
		glog.Error(err.Error())
		return err
	}

	pod.ProviderID = v.UUID
	v.Pods[podId] = pod
	return nil
}

func (n *Node) DeleteVM(vnodeId string) error {
	vnode, exist := n.VMs[vnodeId]
	if !exist {
		err := fmt.Errorf("Node[%s] deletes VM[%s] failed: VM is not found.", n.Name, vnodeId)
		glog.Error(err.Error())
		return err
	}

	vnode.ProviderID = emptyProvider
	delete(n.VMs, vnodeId)

	return nil
}

func (n *Node) AddVM(vnode *VNode) error {
	vnodeId := vnode.UUID

	if _, exist := n.VMs[vnodeId]; exist {
		err := fmt.Errorf("Node[%s] add VM[%s] failed: VM already on the node.", n.Name, vnode.Name)
		glog.Error(err.Error())
		return err
	}

	vnode.ProviderID = n.UUID
	n.VMs[vnodeId] = vnode
	return nil
}

func (c *Container) SetCapacity(cpu, memory float64) error {
	if cpu > 0.0 {
		c.CPU.Capacity = cpu
	}

	if memory > 0.0 {
		c.Memory.Capacity = memory
	}
	return nil
}
