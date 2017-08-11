"""This tool is used to generate topology.conf file for the virtual Cluster.

  Input: The input is the definition of the entities and their amount in the cluster.
        Each line of the input is a json string; Example of the input file:
----the content of an example input file begins----
container, {"kind": "container", "key": "container-cpu", "cpu":[800, 500, 600], "mem": [100, 50, 60], "qps": [120, 50]}
container, {"kind": "container", "key": "container-mem", "cpu":[100, 50, 60], "mem": [800, 500, 600], "qps": [120, 50]}
container, {"kind": "container", "key": "container-cpu-mem", "cpu":[800, 500, 600], "mem": [800, 500, 600], "qps": [120, 90]}
container, {"kind": "container", "key": "container-log", "cpu":[100, 50, 60], "mem": [100, 50, 60], "qps": [120, 1]}

pod, {"kind": "pod", "key": "pod1", "containers":["container-cpu"]}
pod, {"kind": "pod", "key": "pod2", "containers":["container-mem"]}
pod, {"kind": "pod", "key": "pod3", "containers":["container-cpu-mem", "container-log"]}

vnode, {"kind": "vnode", "key": "vnode1", "cpu": 5200, "mem": 8192, "pods":["pod1", "pod1", "pod2"]}
vnode, {"kind": "vnode", "key": "vnode2", "cpu": 5200, "mem": 8192, "pods":["pod2", "pod3", "pod3"]}

node, {"kind": "node", "key": "node1", "cpu": 10400, "mem": 16384, "vnodes": ["vnode1", "vnode1", "vnode2"], "num": 5}
node, {"kind": "node", "key": "node2", "cpu": 10400, "mem": 16384, "vnodes": ["vnode1", "vnode2"], "num": 5}

-----input example ends----


  Output: the topology.conf file of the virtual cluster; the output file format is similar to the input format:

----- output example begins ------
#1. define containers, container format:
# container, <containerId>, <limitCPU>, <usedCPU>, <reqCPU>, <limityMem>, <usedMem>, <reqMem>, <limitQPS>, <usedQPS>;
container, containerC, 300, 180, 100, 400, 350, 250, 100, 80

#2. define Pod, pod format:
# pod, <podId>, <cotainerId1>, <containerId2>
pod, pod-3, containerC


#3. define virtual machine (vnode), vnode format:
# vnode, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <podId1>, <podId2>, ...
vnode, vnode-2, 5200, 8192, 192.168.1.3, pod-3
vnode, vnode-3, 5200, 8192, 192.168.1.4

#4. define the physical machine (node), node format:
# node, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <vnodeId1>, <vnodeId2>, ...
node, node-2, 10400, 16384, 200.0.0.2, vnode-2, vnode-3

#5. define service, service format:
# service, <serviceId>, <podId1>, <podId2>, ...
service, service-1, pod-3
----- output example ends ------

NOTE: currently, this tool will not generate entity of kind service.

"""

import sys
import json


class Container:
    def __init__(self):
        self.kind = "container"
        self.uuid = ""

        self.cpu_limit = 0
        self.cpu_used = 0
        self.cpu_req = 0

        self.mem_limit = 0
        self.mem_used = 0
        self.mem_req = 0

        self.qps_limit = 0
        self.qps_used = 0

        self.index = 0
        return

    def toString(self):
        head = "%s, %s" % (self.kind, self.uuid)
        cpu = "%s, %s, %s" % (self.cpu_limit, self.cpu_used, self.cpu_req)
        mem = "%s, %s, %s" % (self.mem_limit, self.mem_used, self.mem_req)
        qps = "%s, %s" % (self.qps_limit, self.qps_used)

        line = "%s, %s, %s, %s" % (head, cpu, mem, qps)
        return line

    def assign(self, entity):
        if entity['kind'] != 'container':
            print("wrong entity type: %s" %(entity))
            return -1

        self.index += 1
        self.uuid = "%s-%s" % (entity['key'], self.index)

        cpu = entity['cpu']
        self.cpu_limit = cpu[0]
        self.cpu_used = cpu[1]
        self.cpu_req = cpu[2]

        mem = entity['mem']
        self.mem_limit = mem[0]
        self.mem_used = mem[1]
        self.mem_req = mem[2]

        qps = entity['qps']
        self.qps_limit = qps[0]
        self.qps_used = qps[1]
        return 0


class Pod:
    def __init__(self):
        self.kind = "pod"
        self.uuid = ""
        self.containers = []

        self.index = 0
        return

    def toString(self):
        head = "%s, %s" % (self.kind, self.uuid)
        others = ", ".join(self.containers)

        line = "%s, %s" % (head, others)
        return line

    def assign(self, entity, containers):
        kind = entity['kind']
        if kind != 'pod':
            print("wrong entity type, not a pod: %s" % (entity))
            return -1

        self.index += 1
        self.uuid = "%s-%s" % (entity['key'], self.index)
        # self.containers = entity['containers']
        self.containers = containers
        return 0


class VNode:
    def __init__(self, init_ip):
        """ init_ip = [x1, x2, x3, x4], x1, x2, x3, x4 are all integers. """
        self.kind = "vnode"
        self.uuid = ""
        self.cpu = 0
        self.mem = 0
        self.ip = [0, 0, 0, 0]
        self.pods = []

        self.index = 0
        self.init_ip = init_ip #[10, 1, 1, 1]
        return

    def toString(self):
        head = "%s, %s" % (self.kind, self.uuid)
        content = "%s, %s" % (self.cpu, self.mem)
        ip = "%s.%s.%s.%s" % (self.ip[0], self.ip[1], self.ip[2], self.ip[3])
        pods = ", ".join(self.pods)
        line = "%s, %s, %s, %s" % (head, content, ip, pods)

        return line

    def assign(self, entity, pods):
        kind = entity['kind']
        if kind != 'vnode':
            print("wrong entity type, not a vnode: %s" % (entity))
            return -1

        self.index += 1
        self.uuid = "%s-%s" % (entity['key'], self.index)

        self.cpu = entity['cpu']
        self.mem = entity['mem']

        self.ip = next_ip(self.init_ip)
        self.init_ip = self.ip

        self.pods = pods
        return 0


class Node:
    def __init__(self, init_ip):
        self.kind = "node"
        self.uuid = ""
        self.cpu = 0
        self.mem = 0
        self.ip = [0, 0, 0, 0]
        self.vnodes = []

        self.index = 0
        self.init_ip = init_ip
        return

    def toString(self):
        head = "%s, %s" % (self.kind, self.uuid)
        content = "%s, %s" % (self.cpu, self.mem)
        ip = "%s.%s.%s.%s" % (self.ip[0], self.ip[1], self.ip[2], self.ip[3])
        vnodes = ", ".join(self.vnodes)
        line = "%s, %s, %s, %s" % (head, content, ip, vnodes)
        return line

    def assign(self, entity, vnodes):
        kind = entity['kind']
        if kind != 'node':
            print("wrong entity type, not a node: %s" % (entity))
            return -1

        self.index += 1
        self.uuid = "%s-%s" % (entity['key'], self.index)

        self.cpu = entity['cpu']
        self.mem = entity['mem']

        self.ip = next_ip(self.init_ip)
        self.init_ip = self.ip

        self.vnodes = vnodes
        return 0


class Service:
    def __init__(self):
        self.kind = "service"
        self.uuid = ""
        self.pods = []

        self.index = 0
        return

    def toString(self):
        head = "%s, %s" % (self.kind, self.uuid)
        pods = ", ".join(self.pods)
        line = "%s, %s" % (head, pods)
        return line

    def assign(self, entity, pods):
        kind = entity['kind']
        if kind != 'service':
            print("wrong entity type, not a service: %s" % (entity))
            return -1

        self.index += 1
        self.uuid = "%s-%s" % (entity['key'], self.index)

        self.pods = pods
        return


class MyTemplate:
    def __init__(self):
        self.containers = {}
        self.pods = {}
        self.vnodes = {}
        self.nodes = {}
        self.services = {}
        self.kinds = {'container': 1, 
                      'pod': 1,
                      'vnode': 1, 
                      'node': 1, 
                      'service':1}
        return

    def add(self, entity):
        if not isinstance(entity, type({})):
            return -1

        kind = entity.get('kind', None)
        if kind is None:
            return -2

        if not self.kinds.has_key(kind):
            return -3

        key = entity['key']
        if kind == 'container':
            self.containers[key] = entity
        elif kind == 'pod':
            self.pods[key] = entity
        elif kind == 'vnode':
            self.vnodes[key] = entity
        elif kind == 'node':
            self.nodes[key] = entity
        elif kind == 'service':
            self.services[key] = entity
        else:
            msg = "unknown kind:%s: %s" % (kind, entity)
            print(msg)
            return -4

        return 0

    def printContent(self):
        print(self.containers)
        print(self.pods)
        print(self.vnodes)
        print(self.nodes)
        print(self.services)
        return


def next_ip(current_ip):
    if len(current_ip) < 4:
        print("current_ip is illegal %s" % (current_ip))
        return None
    i = 3
    while i >= 0:
        current_ip[i] += 1
        if current_ip[i] <= 254:
            break
        current_ip[i] = 2
        i = i - 1

    if i < 0:
        print("cannot generate next valid ip from: %s", current_ip)
        return None

    return current_ip


class Context:
    def __init__(self, vmIP, pmIP):
        self.containerModel = Container()
        self.podModel = Pod()
        self.vnodeModel = VNode(vmIP)
        self.nodeModel = Node(pmIP)
        self.serviceModel = Service()

        return


def gen_topology(fname_in, fname_out):
    fin = open(fname_in, 'r')
    fout = open(fname_out, 'w')

    template = MyTemplate()

    for line in fin:
        line = line.strip()
        if len(line) < 1 or line[0] == '#':
            continue

        idx = line.find(",")
        if idx < 2:
            print("wrong line format %d :%s" %(idx, line))
            continue
        jstr = line[idx+1:]
        entity = json.loads(jstr)
        template.add(entity)

    template.printContent()
    gen_entities(template, fout)
    fout.close()
    fin.close()
    return 0

# ----------
def gen_container(context, template, ckey, fout):
    containerEntity = template.containers[ckey]

    context.containerModel.assign(containerEntity)
    line = context.containerModel.toString()
    fout.write(line + '\n')
    return context.containerModel.uuid


def gen_pod(context, template, podKey, fout):
    podEntity = template.pods[podKey]
    containerIDs = []

    for containerKey in podEntity['containers']:
        containerId = gen_container(context, template, containerKey, fout)
        containerIDs.append(containerId)

    context.podModel.assign(podEntity, containerIDs)
    line = context.podModel.toString()
    fout.write(line + '\n')
    return context.podModel.uuid


def gen_vnode(context, template, vnodeKey, fout):
    vnodeEntity = template.vnodes[vnodeKey]
    podIDs = []

    for podKey in vnodeEntity['pods']:
        podId = gen_pod(context, template, podKey, fout)
        podIDs.append(podId)

    context.vnodeModel.assign(vnodeEntity, podIDs)
    line = context.vnodeModel.toString()
    fout.write("# vnode, %s\n" % (context.vnodeModel.uuid))
    fout.write(line + '\n\n')
    return context.vnodeModel.uuid


def gen_entities(template, fout):

    vmIP = [10, 1, 1, 1]
    pmIP = [200, 1, 1, 1]
    context = Context(vmIP, pmIP)

    for nodeEntity in template.nodes.itervalues():
        print("node-%s num %s"% (nodeEntity['key'], nodeEntity['num']))
        for i in range(nodeEntity['num']):
            vnodeIDs = []

            for vnodeKey in nodeEntity['vnodes']:
                vnodeId = gen_vnode(context, template, vnodeKey, fout)
                vnodeIDs.append(vnodeId)

            context.nodeModel.assign(nodeEntity, vnodeIDs)
            line = context.nodeModel.toString()
            fout.write("# node, %s\n" % (context.nodeModel.uuid))
            fout.write(line + '\n\n\n')

    return


def main(args):
    print("%s" % (args))
    if len(args) < 3:
        print("Usage:%s <input> <output>" % (args[0]))
        return -1 

    fname_in = args[1]
    fname_out = args[2]
    gen_topology(fname_in, fname_out)
    return 0 

if __name__ == "__main__":
    sys.exit(main(sys.argv))
