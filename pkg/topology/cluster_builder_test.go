package topology

import (
	"testing"
)

func TestNewClusterBuilder(t *testing.T) {
	fname := "../../conf/topology.conf"

	builder := NewClusterBuilder("clusterId-1", "testCluster", fname)
	if builder == nil {
		t.Errorf("load topology failed: %s", fname)
	}

	cluster, err := builder.GenerateCluster()
	if err != nil {
		t.Errorf("failed to generate cluster: %v", err)
	}

	_, err = cluster.GenerateDTOs()
	if err != nil {
		t.Errorf("failed to generate DTOs: %v", err)
	}
}
