package target

const (
	KindApp        = "application"
	KindContainer  = "container"
	KindPod        = "pod"
	KindVirtualApp = "service"
	KindNode       = "host"
	KindCluster = "cluster"
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

	CPU    Resource
	Memory Resource
	Transaction float64
}

type Container struct {
	ObjectMeta

	CPU    Resource
	Memory Resource

	App    *Application
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
	Nodes []*HostNode
	Services []*VirtualApp
}

func NewContainer(name, id string) *Container {
	return &Container{
		ObjectMeta.Kind: KindContainer,
		ObjectMeta.Name: name,
		ObjectMeta.UUID: id,
	}
}

func NewPod(name, id string) *Pod {
	return &Pod{
		ObjectMeta.Kind: KindPod,
		ObjectMeta.Name: name,
		ObjectMeta.UUID: id,
	}
}

func NewHostNode(name, id string) *HostNode {
	return &HostNode{
		ObjectMeta.Kind: KindNode,
		ObjectMeta.Name: name,
		ObjectMeta.UUID: id,
	}
}

func NewApplication(name, id string) *Application {
	return &Application{
		ObjectMeta.Kind: KindApp,
		ObjectMeta.Name: name,
		ObjectMeta.UUID: id,
	}
}

func NewVirtualApp(name, id string) *VirtualApp {
	return &VirtualApp{
		ObjectMeta.Kind: KindVirtualApp,
		ObjectMeta.Name: name,
		ObjectMeta.UUID: id,
	}
}

func NewCluster(name, id string) *Cluster {
	return &Cluster{
		ObjectMeta.Kind: KindCluster,
		ObjectMeta.Name: name,
		ObjectMeta.UUID: id,
	}
}
