# VirtualCluster
1. Generate a virtual cluster with physical machine, virtual machine, pod, container, and servies.
2. Build a supply chain of *physical Node --> virtual Node --> pod --> container --> application --> virtualApplication*.
3. Can execute actions from OpsMgr, and change the topology of this virtualCluster.

# Overview
<div >
<img width="700" src="https://github.com/songbinliu/containerChain/blob/master/conf/supplyChain.png">
</div>

# Commodities between layers
<div>
<img width="500" src="https://github.com/songbinliu/containerChain/blob/master/conf/commodity.png">
</div>

How to decide the amount of commodity bought and sold?

|SE type| vCPU/vMem | CommoditySold | CommodityBought |
|-|-|-|-|
| Application | - | - | Used=Container.Sold.Used |
|Container | Limit/Request/Used | Capacity=Limit (if no limit, then pod.Capacity) <br/> Used=*Monitored-Container* | Used=Container.Sold.Used|
|Pod | Capacity/Used | Capacity=VM.Capacity  <br/> Used=sum.Container.Bought.Used | Used=sum.Container.Bought.Used |
|VM | Capacity/Used | Capacity=Capacity <br/> Used=*Monitored-VM* | Used=VM.Capacity|
|PM | Capacity/Used | Capacity=Capacity<br/> Used=*Monitored-PM*| -|

1.*Monitored-Container* : monitored resource usage of container;<br/>
2.*Monitored-VM*: monitored resource usage of VM (= sum.Pod.Bought.Used + overhead1);<br/>
3.*Monitored-PM*: monitored resource usage of PM (= sum.Monitored-VM + overhead2);<br/>
4.*Container.Limit and Container.Request* are read from Container settings.<br/>


# Supported Actions
|SE type| Move | Resize|
|-|-|-|
|ContainerPod| Yes | No |
|Container | No | Yes |
| VirtualMachine |Yes | WIP|

 (*WIP* = work in progress.)

# Run it

```console
#1. get source code
go get github.com/turbonomic/virtualCluster

#2. compile it
cd $GOPATH/src/github.com/turbonomic/virtualCluster
make build

#3. run it
turbo=./conf/turbo.json
topology=./conf/topology.conf
target=./conf/target.json
./_output/vCluster --topologyConf $topology --turboConf $turbo --targetConf $target --logtostderr --v 3 

Note: in case of updating dependency, run glide before compiling it:
## update the version of turbo-go-sdk in glide.yaml first if necessary
glide update --strip-vendor
```

**turbo** is a json file about the settings of the OpsMgr, [example](https://github.com/turbonomic/virtualCluster/blob/master/conf/turbo.json);

**target** is a json file about settings of generated cluster for OpsMgr, [example](https://github.com/turbonomic/virtualCluster/blob/master/conf/target.json);

**topology** is the configuration file about the virtual cluster to be generated, [example](https://github.com/turbonomic/virtualCluster/blob/master/conf/topology.conf).

# Topologies
Different topologies will trigger different actions from OpsMgr.

## Resize Up containers
```
#1. define containers, container format:
# container, <containerId>, <limitCPU>, <usedCPU>, <reqCPU>, <limityMem>, <usedMem>, <reqMem>, <limitQPS>, <usedQPS>;
container, containerC, 1000, 900, 500, 1624, 224, 250, 100, 10
container, containerD, 2900, 900, 500, 1024, 950, 250, 100, 20

#2. define Pod, pod format:
# pod, <podId>, <cotainerId1>, <containerId2>
pod, pod-3, containerC
pod, pod-4, containerD

#3. define virtual machine (vnode), vnode format:
# vnode, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <podId1>, <podId2>, ...
vnode, vnode-3, 5200, 4096, 192.168.1.4, pod-4, pod-3

#4. define the physical machine (node), node format:
# node, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <vnodeId1>, <vnodeId2>, ...
node, node-3, 10400, 16384, 200.0.0.2, vnode-3
```
In [this topology](https://github.com/turbonomic/virtualCluster/blob/3a2153cb3eef21fc6cdb20945eee5d971e671b36/conf/resize.up.container.topology.conf#L13), the CPU utilization of `containerC` is high, so an action will be triggered to increase the CPU capacity of `containerC`; and another action to increase the Memory capacity for `containerD`.


## Move Pods to an Idle VM
```
#1. define containers, container format:
# container, <containerId>, <limitCPU>, <usedCPU>, <reqCPU>, <limityMem>, <usedMem>, <reqMem>, <limitQPS>, <usedQPS>;
container, containerC, 2000, 950, 500, 2048, 1024, 250, 100, 10
container, containerD, 2000, 850, 500, 2048, 800, 250, 100, 20

#2. define Pod, pod format:
# pod, <podId>, <cotainerId1>, <containerId2>
pod, pod-3, containerC
pod, pod-4, containerD
pod, pod-5, containerC
pod, pod-6, containerD

#3. define virtual machine (vnode), vnode format:
# vnode, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <podId1>, <podId2>, ...
vnode, vnode-2, 4200, 4096, 192.168.1.3
vnode, vnode-3, 4200, 4096, 192.168.1.4, pod-4, pod-3, pod-5, pod-6

#4. define the physical machine (node), node format:
# node, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <vnodeId1>, <vnodeId2>, ...
node, node-3, 10400, 16384, 200.0.0.2, vnode-2, vnode-3
```
In [this topology](https://github.com/turbonomic/virtualCluster/blob/3a2153cb3eef21fc6cdb20945eee5d971e671b36/conf/move.pod.topology.conf#L13), `vnode-3` is hosting four pods, and is highly utilized; on the other hand, `vnode-2` is idle. So actions will be triggered to move some pods from `vnode-3` to `vnode-2`.

## Provision VM
```
#1. define containers, container format:
# container, <containerId>, <limitCPU>, <usedCPU>, <reqCPU>, <limityMem>, <usedMem>, <reqMem>, <limitQPS>, <usedQPS>;
container, containerA, 700, 700, 700, 700, 700, 100, 100, 0
container, containerB, 2300, 2300, 2300, 2300, 2300, 100, 100, 0

#2. define Pod, pod format:
# pod, <podId>, <cotainerId1>, <containerId2>
pod, pod-1, containerB
pod, pod-1-2, containerA
pod, pod-2, containerA

#3. define virtual machine (vnode), vnode format:
# vnode, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <podId1>, <podId2>, ...
vnode, vnode-1, 3000, 3000, 192.168.1.2, pod-1, pod-1-2
vnode, vnode-2, 1100, 1100, 192.168.1.3, pod-2

#4. define the physical machine (node), node format:
# node, <nodeId>, <cpu_capacity>, <mem_capacity>, <IP>, <vnodeId1>, <vnodeId2>, ...
node, node-1, 3200, 3200, 200.0.0.1, vnode-1
node, node-2, 2100, 2100, 200.0.0.2, vnode-2

#4. define service, service format:
# service, <serviceId>, <podId1>, <podId2>, ...
service, service-1, pod-1, pod-2

```
In [this topology](./conf/provision.vm.topology.conf), `vnode-1` is hosting two pods, and is highly utilized; on the other hand, `vnode-2` does not have enough resource available for the pods hosted by `vnode-1` to move to. In such scenario, since `node-2` has sufficient resource available,  actions will be triggered to provision a VM in `node-2` so that one pod in `vnode-1` can move to the new VM to reduce its critical utilization.
