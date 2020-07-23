package discovery

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
)

// Configuration Parameters to connect to an example Target
type TargetConf struct {
	Address         string
	Username        string
	Password        string
	ProbeCategory   string
	TargetType      string
	ProbeUICategory string
}

// Create a new ExampleClientConf from file. Other fields have default values and can be overridden.
func NewTargetConf(path string) (*TargetConf, error) {

	glog.Infof("[TargetConf] Read configuration from %s\n", path)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Errorf("failed to read file:%v", err.Error())
		return nil, err
	}

	var config TargetConf
	err = json.Unmarshal(file, &config)

	if err != nil {
		msg := fmt.Sprintf("Unmarshall error :%v\n", err)
		glog.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	glog.V(2).Infof("Results: %+v\n", config)

	return &config, nil
}
