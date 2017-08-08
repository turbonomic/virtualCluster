package target

const (
	KindApp        = "application"
	KindContainer  = "container"
	KindPod        = "pod"
	KindVirtualApp = "service"
	KindVNode       = "vhost"
	KindNode       = "host"
	KindCluster    = "cluster"
)

type ObjectMeta struct {
	Name string
	UUID string
	Kind string
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

	CPU    Resource
	Memory Resource

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

	CPU       Resource
	Memory    Resource

	ClusterId string
	IP        string

	//a map for easy of move/deletion, key=pod.UUID
	Pods      map[string]*Pod
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
