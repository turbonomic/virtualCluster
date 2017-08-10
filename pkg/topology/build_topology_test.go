package topology

import (
	"testing"
)

func TestTargetTopology_LoadTopology(t *testing.T) {
	fname := "../../conf/topology.conf"

	topo := NewTargetTopology("testCluster")

	err := topo.LoadTopology(fname)
	if err != nil {
		t.Errorf("load topology test failed. %v", err)
	}
}
