# containerChain
build a supply chain of *physical Node --> virtual Node --> pod --> container --> application --> virtualApplication*.

# Overview
<div >
<img width="700" src="https://github.com/songbinliu/containerChain/blob/master/conf/supplyChain.png">
</div>

# Commodities between layers
<div>
<img width="500" src="https://github.com/songbinliu/containerChain/blob/master/conf/commodity.png">
</div>

# Supported Actions
|SE type| Move | Resize|
|-|-|-|
|ContainerPod| Yes | No |
|Container | No | WIP |
| VirtualMachine |Yes | WIP|

note *WIP* = work in progress.

# Run it

```bash
#1. get source code
go get github.com/songbinliu/containerChain

#2. compile it
cd $GOPATH/src/github.com/songbinliu/containerChain
make build

#3. run it
./_output/containerChain --logtostderr --v 3 --topologyConf ./conf/simple.topology.conf --turboConf ./conf/turbo.json --targetConf ./conf/target.json
```
