package topology

import (
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/virtualCluster/pkg/target"
	"os"
	"strconv"
	"strings"
	"regexp"
)

const (
	defaultQPSLimit = float64(120)
)

type containerTemplate struct {
	Key    string
	CPU    target.Resource
	Memory target.Resource
	ReqCPU float64
	ReqMem float64

	QPS target.Resource
	ResponseTime *target.Resource
}

type podTemplate struct {
	Key        string
	Containers []string
}

// virtual machine
type vnodeTemplate struct {
	Key string

	CPU    float64
	Memory float64
	IP     string
	Pods   []string
}

// physical machine
type nodeTemplate struct {
	Key string

	CPU    float64
	Memory float64
	IP     string
	VMs    []string
}

type serviceTemplate struct {
	Key  string
	Pods []string
}

type TargetTopology struct {
	ClusterId string

	// containerTemplate map
	ContainerTemplateMap map[string]*containerTemplate

	// podTemplate map
	PodTemplateMap map[string]*podTemplate

	//nodeTemplate map
	VNodeTemplateMap map[string]*vnodeTemplate

	//physicalMachine map
	NodeTemplateMap map[string]*nodeTemplate

	//serviceTemplate amp
	ServiceTemplateMap map[string]*serviceTemplate
}

func NewTargetTopology(clusterId string) *TargetTopology {
	topo := &TargetTopology{
		ClusterId:            clusterId,
		ContainerTemplateMap: make(map[string]*containerTemplate),
		PodTemplateMap:       make(map[string]*podTemplate),
		VNodeTemplateMap:     make(map[string]*vnodeTemplate),
		NodeTemplateMap:      make(map[string]*nodeTemplate),
		ServiceTemplateMap:   make(map[string]*serviceTemplate),
	}

	return topo
}

// load containerTemplate from a line
// fields: containerName, req_cpu, used_cpu, req_memory, used_mem [, qpsLimit, qpsUsed [, responseTimeCap, responseTimeUsed ] ]
func loadContainer(t *TargetTopology, input *InputLine) error {
	if _, exist := t.ContainerTemplateMap[input.key]; exist {
		return fmt.Errorf("container[%s] already exists.", input.key)
	}

	// CPU amount
	limitCPU := input.getFloat()
	usedCPU := input.getFloat()
	reqCPU := input.getFloat()

	limitMem := input.getFloat()
	usedMem := input.getFloat()
	reqMem := input.getFloat()

	// QPS amount
	limitQPS := defaultQPSLimit
	usedQPS := 0.0
	if input.RemainingFieldCount() >= 1 {
		limitQPS = input.getFloat()
		usedQPS = input.getFloat()
	}

	var responseTimeCommodity *target.Resource
	if input.RemainingFieldCount() >= 1 {
		responseTimeCommodity = &target.Resource{}
		responseTimeCommodity.Capacity = input.getFloat()
		responseTimeCommodity.Used = input.getFloat()
	}

	container := &containerTemplate{
		Key: input.key,
		CPU: target.Resource{
			Capacity: limitCPU,
			Used:     usedCPU,
		},
		ReqCPU: reqCPU,

		// change the unit of Memory from MB to KB
		Memory: target.Resource{
			Capacity: limitMem * 1024.0,
			Used:     usedMem * 1024.0,
		},
		ReqMem: reqMem * 1024.0,

		QPS: target.Resource{
			Capacity: limitQPS,
			Used:     usedQPS,
		},
		ResponseTime: responseTimeCommodity,
	}

	if input.err == nil {
		t.ContainerTemplateMap[input.key] = container
		glog.V(4).Infof("[container] %+v", container)
	}
	return input.err
}

// load podTemplate from a line
// pod, key, container1, container2, ...
func loadPod(t *TargetTopology, input *InputLine) error {
	if _, exist := t.PodTemplateMap[input.key]; exist {
		err := fmt.Errorf("Pod[%s] already exists", input.key)
		glog.Error(err.Error())
		return err
	}

	if input.RemainingFieldCount() < 1 {
		return fmt.Errorf("missing container list in pod declaration")
	}
	pod := &podTemplate{
		Key:        input.key,
		Containers: input.GetRestOfFields(),
	}

	t.PodTemplateMap[input.key] = pod
	glog.V(4).Infof("[pod] %+v", pod)
	return nil
}

// load vnodeTemplate from a line
// vnode.key, cpu, memory, IP, pod1, pod2, ...
func loadVNode(t *TargetTopology, input *InputLine) error {
	if _, exist := t.VNodeTemplateMap[input.key]; exist {
		err := fmt.Errorf("vnode [%s] already exists", input.key)
		glog.Error(err.Error())
		return err
	}

	cpu := input.getFloat()
	mem := input.getFloat()
	ip := input.getString()

	if input.RemainingFieldCount() < 1 {
		return fmt.Errorf("missing pod list in vnode declaration")
	}

	vnode := &vnodeTemplate{
		Key:    input.key,
		CPU:    cpu,
		Memory: mem * 1024.0,
		IP:     ip,
		Pods:   input.GetRestOfFields(),
	}

	t.VNodeTemplateMap[input.key] = vnode
	glog.V(4).Infof("[vnode] %+v", vnode)
	return nil
}

// load nodeTemplate from a line
// node.key, cpu, memory, IP, vnode1, vnode2, ...
func loadNode(t *TargetTopology, input *InputLine) error {
	if _, exist := t.NodeTemplateMap[input.key]; exist {
		err := fmt.Errorf("node [%s] already exists", input.key)
		glog.Error(err.Error())
		return err
	}

	cpu := input.getFloat()
	mem := input.getFloat()
	ip := input.getString()

	if input.RemainingFieldCount() < 1 {
		return fmt.Errorf("missing vnode list in node declaration")
	}

	node := &nodeTemplate{
		Key:    input.key,
		CPU:    cpu,
		Memory: mem * 1024,
		IP:     ip,
		VMs:    input.GetRestOfFields(),
	}

	t.NodeTemplateMap[input.key] = node
	glog.V(4).Infof("[node] %+v", node)
	return nil
}

// load serviceTemplate from a line
// service-key, pod1, pod2, ...
func loadService(t *TargetTopology, input *InputLine) error {
	if _, exist := t.ServiceTemplateMap[input.key]; exist {
		err := fmt.Errorf("service[%s] already exists", input.key)
		glog.Error(err.Error())
		return err
	}

	if input.RemainingFieldCount() < 1 {
		return fmt.Errorf("missing pod list in service declaration")
	}

	service := &serviceTemplate{
		Key:  input.key,
		Pods: input.GetRestOfFields(),
	}

	t.ServiceTemplateMap[input.key] = service
	glog.V(4).Infof("[service] %+v", service)
	return nil
}

type InputLine struct {
	err error
	line string			// original line
	key string
	fields []string
	command string
	fieldNum int
}

/*
 * Prepare a line from a topology definition file for parsing
 */
var commentPattern = regexp.MustCompile("#.*")

func makeInputLine (line string) (*InputLine, error) {
	var err error
	il := InputLine{line: line}
	commentsRemoved := commentPattern.ReplaceAllString(line, "")
	if len(commentsRemoved) == 0 {
		// A line with only a comment has a pseudo entity type of "comment"
		commentsRemoved = "comment, key"
	}
	for i, field := range strings.Split(commentsRemoved, ",") {
		trimmed := strings.TrimSpace(field)
		if len(trimmed) == 0 {
			err = fmt.Errorf("field %d is empty", i + 1)
			break
		}
		il.fields = append(il.fields, trimmed)
	}
	il.command = il.getString()
	il.key = il.getString()
	if il.err != nil && err == nil {
		err = fmt.Errorf("missing key field")
	}
	return &il, err
}

func (l *InputLine) getString() string {
	value := ""
	if l.err == nil {
		if l.fieldNum >= len(l.fields) {
			l.err = fmt.Errorf("input line '%s' has insufficient fields", l.line)
		} else {
			value = l.fields[l.fieldNum]
			l.fieldNum++
		}
	}
	return value
}

func (l *InputLine) getFloat() float64 {
	value := 0.0
	var err error
	s := l.getString()
	if l.err == nil {
		if value, err = strconv.ParseFloat(s, 64); err != nil {
			l.err = fmt.Errorf("invalid float value '%s' at field %d", s, l.fieldNum)
		}
	}
	return value
}
func (l *InputLine) RemainingFieldCount() int {
	return len(l.fields) - l.fieldNum
}

func (l *InputLine) GetRestOfFields() []string {
	strlist := l.fields[l.fieldNum:]
	l.fieldNum = len(l.fields)
	return strlist
}

type HandlerFunction func(*TargetTopology, *InputLine)error
func noop(_ *TargetTopology, _ *InputLine) error {
	return nil
}
var loadHandlers = map[string]HandlerFunction{
	"container": loadContainer,
	"pod": loadPod,
	"vnode": loadVNode,
	"node": loadNode,
	"service": loadService,
	"comment": noop,
}
func (t *TargetTopology) parseLine(lineNum int, input *InputLine) error {
	var err error
	handler := loadHandlers[input.command]
	if handler != nil {
		glog.V(4).Infof("begin to build a container [%d]: %s", lineNum, input.line)
		err = handler(t, input)
		if err == nil && input.RemainingFieldCount() > 0 {
			err = fmt.Errorf("line %d has unused fields", lineNum)
		}
	} else {
		err = fmt.Errorf("invalid EntityType[%s]", input.command)
	}

	return err
}

func (t *TargetTopology) CheckTemplateEmpty() error {
	if len(t.PodTemplateMap) < 1 {
		err := fmt.Errorf("podTemplate is empty.")
		glog.Error(err.Error())
		return err
	}

	if len(t.ContainerTemplateMap) < 1 {
		err := fmt.Errorf("containerTemplate is empty.")
		glog.Error(err.Error())
		return err
	}

	if len(t.VNodeTemplateMap) < 1 {
		err := fmt.Errorf("vnodeTemplate is empty.")
		glog.Error(err.Error())
		return err
	}

	if len(t.NodeTemplateMap) < 1 {
		err := fmt.Errorf("nodeTemplate is empty.")
		glog.Error(err.Error())
		return err
	}

	if len(t.ServiceTemplateMap) < 1 {
		glog.Warningf("serviceTemplateMap is empty.")
	}

	return nil
}

func (t *TargetTopology) PrintTemplateInfo() {
	glog.V(1).Infof("containerTemplate.num=%d", len(t.ContainerTemplateMap))
	glog.V(1).Infof("podTemplate.num=%d", len(t.PodTemplateMap))
	glog.V(1).Infof("vnodeTemplate.num=%d", len(t.VNodeTemplateMap))
	glog.V(1).Infof("nodeTemplate.num=%d", len(t.NodeTemplateMap))
	glog.V(1).Infof("serviceTemplate.num=%d", len(t.ServiceTemplateMap))
}

func (t *TargetTopology) LoadTopology(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		glog.Errorf("failed to open file[%s] for read: %v", fname, err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum += 1
		input, err := makeInputLine(scanner.Text())
		if err == nil {
			err = t.parseLine(lineNum, input)
		}
		if err != nil {
			glog.Errorf("parse [%s/%d] line[%s] failed: %v", fname, lineNum, input.line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		glog.Errorf("error while reading file[%s]: %v", fname, err)
		return err
	}

	if err := t.CheckTemplateEmpty(); err != nil {
		err := fmt.Errorf("Template checked failed: %v", err)
		glog.Error(err.Error())
		return err
	}

	t.PrintTemplateInfo()

	return nil
}
