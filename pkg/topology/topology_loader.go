package topology

import (
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"github.com/songbinliu/containerChain/pkg/target"
	"os"
	"strconv"
	"strings"
)

type containerTemplate struct {
	Key    string
	CPU    target.Resource
	Memory target.Resource
	ReqCPU float64
	ReqMem float64
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
//fields: containerName, req_cpu, used_cpu, req_memory, used_mem
func (t *TargetTopology) loadContainer(fields []string) error {
	expectNumFields := 7
	if len(fields) != expectNumFields {
		return fmt.Errorf("fields num mismatch [%d Vs. %d]", len(fields), expectNumFields)
	}
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
		if len(fields[i]) < 1 {
			return fmt.Errorf("field[%d] of fields-%v is empty", i+1, fields)
		}
	}

	i := 0
	key := fields[i]
	i++
	if _, exist := t.ContainerTemplateMap[key]; exist {
		return fmt.Errorf("container[%s] already exists.", key)
	}

	limitCPU, err := strconv.ParseFloat(fields[i], 64)
	i++
	if err != nil {
		return fmt.Errorf("limit_cpu field-%d-[%s] should be a float number.", i, fields[i-1])
	}

	usedCPU, err := strconv.ParseFloat(fields[i], 64)
	i++
	if err != nil {
		return fmt.Errorf("used_cpu field-%d-[%s] should be a float number.", i, fields[i-1])
	}

	reqCPU, err := strconv.ParseFloat(fields[i], 64)
	i++
	if err != nil {
		return fmt.Errorf("req_cpu field-%d-[%s] should be a float number.", i, fields[i-1])
	}

	limitMem, err := strconv.ParseFloat(fields[i], 64)
	i++
	if err != nil {
		return fmt.Errorf("limit_mem field-%d-[%s] should be a float number.", i, fields[i-1])
	}

	usedMem, err := strconv.ParseFloat(fields[i], 64)
	i++
	if err != nil {
		return fmt.Errorf("used_mem field-%d-[%s] should be a float number.", i, fields[i-1])
	}

	reqMem, err := strconv.ParseFloat(fields[i], 64)
	i++
	if err != nil {
		return fmt.Errorf("req_mem field-%d-[%s] should be a float number.", i, fields[i-1])
	}

	container := &containerTemplate{
		Key: key,
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
	}

	t.ContainerTemplateMap[key] = container
	glog.V(3).Infof("[container] %+v", container)
	return nil
}

// load podTemplate from a line
// pod.key, container1, container2, ...
func (t *TargetTopology) loadPod(fields []string) error {
	expectNumFields := 2
	if len(fields) < expectNumFields {
		return fmt.Errorf("fields too fewer [%d Vs. %d]", len(fields), expectNumFields)
	}

	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
		if len(fields[i]) < 1 {
			return fmt.Errorf("field[%d] of fields-%v is empty", i+1, fields)
		}
	}

	key := fields[0]
	if _, exist := t.PodTemplateMap[key]; exist {
		fmt.Errorf("Pod[%s] already exist.")
	}

	containers := []string{}
	for i := 1; i < len(fields); i++ {
		containers = append(containers, fields[i])
	}

	pod := &podTemplate{
		Key:        key,
		Containers: containers,
	}

	t.PodTemplateMap[key] = pod
	glog.V(4).Infof("[pod] %+v", pod)
	return nil
}

// load vnodeTemplate from a line
// vnode.key, cpu, memory, IP, pod1, pod2, ...
func (t *TargetTopology) loadVNode(fields []string) error {
	expectNumFields := 4
	if len(fields) < expectNumFields {
		return fmt.Errorf("fields too fewer [%d Vs. %d]", len(fields), expectNumFields)
	}

	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
		if len(fields[i]) < 1 {
			return fmt.Errorf("field[%d] of fields-%v is empty", i+1, fields)
		}
	}

	key := fields[0]
	if _, exist := t.VNodeTemplateMap[key]; exist {
		fmt.Errorf("vnode [%s] already exist.")
	}

	cpu, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return fmt.Errorf("convert field-1-cpu[%s] failed: %v", fields[1], err)
	}

	mem, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return fmt.Errorf("conver field-2-mem[%s] failed: %v", fields[2], err)
	}

	ip := fields[3]

	pods := []string{}
	for i := 4; i < len(fields); i++ {
		pods = append(pods, fields[i])
	}

	vnode := &vnodeTemplate{
		Key:    key,
		CPU:    cpu,
		Memory: mem * 1024.0,
		IP:     ip,
		Pods:   pods,
	}

	t.VNodeTemplateMap[key] = vnode
	glog.V(4).Infof("[vnode] %+v", vnode)
	return nil
}

// load nodeTemplate from a line
// node.key, cpu, memory, IP, vnode1, vnode2, ...
func (t *TargetTopology) loadNode(fields []string) error {
	expectNumFields := 4
	if len(fields) < expectNumFields {
		return fmt.Errorf("fields too fewer [%d Vs. %d]", len(fields), expectNumFields)
	}

	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
		if len(fields[i]) < 1 {
			return fmt.Errorf("field[%d] of fields-%v is empty", i+1, fields)
		}
	}

	key := fields[0]
	if _, exist := t.NodeTemplateMap[key]; exist {
		fmt.Errorf("node [%s] already exist.")
	}

	cpu, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return fmt.Errorf("convert field-1-cpu[%s] failed: %v", fields[1], err)
	}

	mem, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return fmt.Errorf("conver field-2-mem[%s] failed: %v", fields[2], err)
	}

	ip := fields[3]

	vms := []string{}
	for i := 4; i < len(fields); i++ {
		vms = append(vms, fields[i])
	}

	node := &nodeTemplate{
		Key:    key,
		CPU:    cpu,
		Memory: mem * 1024,
		IP:     ip,
		VMs:    vms,
	}

	t.NodeTemplateMap[key] = node
	glog.V(4).Infof("[node] %+v", node)
	return nil
}

// load serviceTemplate from a line
// service-key, pod1, pod2, ...
func (t *TargetTopology) loadService(fields []string) error {
	expectNumFields := 2
	if len(fields) < expectNumFields {
		return fmt.Errorf("fields too fewer [%d Vs. %d]", len(fields), expectNumFields)
	}

	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
		if len(fields[i]) < 1 {
			return fmt.Errorf("field[%d] of fields-%v is empty", i+1, fields)
		}
	}

	key := fields[0]
	if _, exist := t.ServiceTemplateMap[key]; exist {
		fmt.Errorf("service[%s] already exist.")
	}

	pods := []string{}
	for i := 1; i < len(fields); i++ {
		pods = append(pods, fields[i])
	}

	service := &serviceTemplate{
		Key:  key,
		Pods: pods,
	}

	t.ServiceTemplateMap[key] = service
	glog.V(4).Infof("[service] %+v", service)
	return nil
}

func (t *TargetTopology) parseLine(lineNum int, line string, fields []string) error {
	entityType := strings.TrimSpace(fields[0])

	var err error
	switch entityType {
	case "container":
		glog.V(2).Infof("begin to build a container [%d]: %s", lineNum, line)
		err = t.loadContainer(fields[1:])
	case "pod":
		glog.V(2).Infof("begin to build a pod [%d]: %s", lineNum, line)
		err = t.loadPod(fields[1:])
	case "vnode":
		glog.V(2).Infof("begin to build a vnode [%d]: %s", lineNum, line)
		err = t.loadVNode(fields[1:])
	case "node":
		glog.V(2).Infof("begin to build a node [%d]: %s", lineNum, line)
		err = t.loadNode(fields[1:])
	case "service":
		glog.V(2).Infof("begin to build a service [%d]: %s", lineNum, line)
		err = t.loadService(fields[1:])
	default:
		err = fmt.Errorf("wrong EntityType[%s]", fields[0])
	}

	if err != nil {
		return fmt.Errorf("build %s failed: %v", entityType, err)
	}

	return nil
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
		line := strings.TrimSpace(scanner.Text())
		lineNum += 1

		if len(line) < 1 || line[0] == '#' {
			glog.V(4).Infof("skip file[%s] line#%d", fname, lineNum)
			continue
		}

		segs := strings.Split(line, ",")
		if len(segs) < 1 {
			glog.V(2).Infof("Invalid file[%s] line#%d", fname, lineNum)
		}

		if err := t.parseLine(lineNum, line, segs); err != nil {
			glog.Errorf("parse [%s/%d] line[%s] failed: %v", fname, lineNum, line, err)
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
