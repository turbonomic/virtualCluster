package target

type ObjectMeta struct {
	Name string
	UUID string
	Kind string
}

type Resource struct {
	Capacity float64
	Used     float64
}

type PhysicalMachine struct {
	ObjectMeta

	CPU    Resource
	Memory Resource

	ClusterID string
	IP        string

	Pods []*Pod
}

type Pod struct {
	ObjectMeta

	CPU    Resource
	Memory Resource

	Containers []*Container
}

type Container struct {
	ObjectMeta

	CPU    Resource
	Memory Resource
}

type Application struct {
	ObjectMeta

	CPU    Resource
	Memory Resource
}

type Cluster struct {
	ObjectMeta
	Nodes []*PhysicalMachine
}
