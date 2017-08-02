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
	AppName string

	CPU    Resource
	Memory Resource

	Containers []*Container
}

type Container struct {
	ObjectMeta

	CPU    Resource
	Memory Resource

	app *Application
}

type Application struct {
	ObjectMeta

	CPU    Resource
	Memory Resource
}

type VApplication struct {
	ObjectMeta

	pods []*Pod
}

type Cluster struct {
	ObjectMeta
	Nodes []*PhysicalMachine
}
