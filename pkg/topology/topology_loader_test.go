package topology

import (
	"testing"
	"fmt"
	"github.com/songbinliu/virtualCluster/pkg/util"
	"bytes"
)

func TestTargetTopology_LoadTopology(t *testing.T) {
	fname := "../../conf/topology.conf"

	topo := NewTargetTopology("testCluster")

	err := topo.LoadTopology(fname)
	if err != nil {
		t.Errorf("load topology test failed. %v", err)
	}
}

var expected []string = []string{
	"parse [../../conf/parse-test.topo/2] line[container, container-2b, 1010.0, 911.0, 512.0, 613.000000, 614.000000, 215.000000, 116.0] failed: input line 'container, container-2b, 1010.0, 911.0, 512.0, 613.000000, 614.000000, 215.000000, 116.0' has insufficient fields",
	"parse [../../conf/parse-test.topo/4] line[pod, pod-2, container-2,,foo] failed: field 4 is empty",
	"parse [../../conf/parse-test.topo/5] line[,pod, pod-x, container-2] failed: field 1 is empty",
	"parse [../../conf/parse-test.topo/7] line[container, container-2d, 1010.0, 911.0, 512.0, 613.000000, 614.000000, 215.000000, 116.0, 17.0, 1, 2, 3, 4] failed: line 7 has unused fields",
	"parse [../../conf/parse-test.topo/15] line[foo, container-3, 1010.0, 911.0, 512.0, 613.000000, 614.000000, 215.000000, 116.1, 17.1, 500, 30] failed: invalid EntityType[foo]",
	"parse [../../conf/parse-test.topo/17] line[container, container-2c, 1010.0, 91z.0, 512.0, 613.000000, 614.000000, 215.000000, 116.0, 17.0] failed: invalid float value '91z.0' at field 4",
	"parse [../../conf/parse-test.topo/19] line[pod, bad-pod] failed: missing container list in pod declaration",
	"parse [../../conf/parse-test.topo/21] line[container, container-1b, 1001.0, 992.0, 503.0, 1624.000000, 225.000000, 256.000000, 107.0, 18.0, 500] failed: input line 'container, container-1b, 1001.0, 992.0, 503.0, 1624.000000, 225.000000, 256.000000, 107.0, 18.0, 500' has insufficient fields",
	"Pod[pod-1] already exists",
	"parse [../../conf/parse-test.topo/24] line[pod, pod-1, container-1] failed: Pod[pod-1] already exists",
	"",
}

func dumplist(strlist []string) string {
	var buf bytes.Buffer
	for _, s := range strlist {
		fmt.Fprint(&buf, "  %s\n", s)
	}
	return buf.String()
}

func TestTargetTopology_LoadTopologyFailures(t *testing.T) {
	fname := "../../conf/parse-test.topo"
	topo := NewTargetTopology("testCluster")

	output, _ := testutil.GetOutputAsList(func() {
		topo.LoadTopology(fname)
	}, 54)
	failed := false
	if len(output) != len(expected) {
		failed = true
	} else {
		for i, s := range expected {
			if s != output[i] {
				failed = true
				break
			}
		}
	}
	if failed {
		t.Errorf("Parse verification failed.  Expected:\n%s\nGot:\n%s", dumplist(expected), dumplist(output))
	}
}
