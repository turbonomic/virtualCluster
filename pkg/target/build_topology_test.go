package target

import (
	"testing"
)

func TestTargetTopology_LoadTopology(t *testing.T) {
	fname := "../conf/small.conf"

	topo := NewTargetTopology("testCluster")

	err := topo.LoadTopology(fname)
	if err != nil {
		t.Errorf("load topology test failed. %v", err)
	}
}
