package target

const (
	KindApp        = "application"
	KindContainer  = "container"
	KindPod        = "pod"
	KindVirtualApp = "service"
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
	AppName string

	CPU    Resource
	Memory Resource

	Containers []*Container
}

type VirtualApp struct {
	ObjectMeta

	Pods []*Pod
}

type HostNode struct {
	ObjectMeta

	CPU    Resource
	Memory Resource

	ClusterID string
	IP        string

	Pods []*Pod
}

type Cluster struct {
	ObjectMeta
	Nodes    []*HostNode
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

func NewHostNode(name, id string) *HostNode {
	return &HostNode{
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
