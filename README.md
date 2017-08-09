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

What is the amount of buy and sell?

|SE type| vCPU/vMem | CommoditySell | CommodityBought |
|-|-|-|-|
| Application | - | - | Used=Container.Used |
|Container | Limit/Request/Used | Capacity=Limit <br/> Used= Request Vs. Used | Used = Limit, Request Vs. Used|
|Pod | Limit/Request/Used | Capacity=VM.Capacity  <br/> Used=Container.(Request Vs. Used) | Used = Container.(Request Vs. Used) |
|VM | Capacity/Used | Capacity=Capacity <br/> Used=Used | Used=Capacity|

*Container.Used* is the monitored usage; others should be calculated based on *Container.Used*.

# Supported Actions
|SE type| Move | Resize|
|-|-|-|
|ContainerPod| Yes | No |
|Container | No | WIP |
| VirtualMachine |Yes | WIP|

 (*WIP* = work in progress.)

# Run it

```console
#1. get source code
go get github.com/songbinliu/containerChain

#2. compile it
cd $GOPATH/src/github.com/songbinliu/containerChain
make build

#3. run it
turbo=./conf/turbo.json
topology=./conf/topology.conf
target=./conf/target.json
./_output/containerChain --topologyConf $topology --turboConf $turbo --targetConf $target --logtostderr --v 3 
```

**turbo** is a json file about the settings of the OpsMgr, [example](https://github.com/songbinliu/containerChain/blob/master/conf/turbo.json);

**target** is a json file about settings of generated cluster for OpsMgr, [example](https://github.com/songbinliu/containerChain/blob/master/conf/target.json);

**topology** is the configuration file about the virtual cluster to be generated, [example](https://github.com/songbinliu/containerChain/blob/master/conf/topology.conf).
