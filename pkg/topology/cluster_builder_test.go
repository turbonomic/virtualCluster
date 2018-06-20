package topology

import (
	"testing"
	"fmt"
	"github.com/turbonomic/virtualCluster/pkg/util"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/stretchr/testify/assert"
)

func generateTestCluster() ([]*proto.EntityDTO, string) {
	fname := testutil.MakeTestPath("conf/topology.conf")
	builder := NewClusterBuilder("clusterId-1", "testCluster", fname)
	if builder == nil {
		return nil, fmt.Sprintf("load topology failed: %s", fname)
	}

	cluster, err := builder.GenerateCluster()
	if err != nil {
		return nil, fmt.Sprintf("failed to generate cluster: %v", err)
	}

	dtoList, err := cluster.GenerateDTOs()
	if err != nil {
		return nil, fmt.Sprintf("failed to generate DTOs: %v", err)
	}
	return dtoList, ""
}

func TestNewClusterBuilder(t *testing.T) {
	_, err := generateTestCluster()
	if (err != "") {
		t.Error(err)
	}
}

func TestResponseTimeCalculation(t *testing.T) {
	dtoList, err := generateTestCluster()
	if (err != "") {
		t.Error(err)
	}
	/* We know that the weighted response time average offered by service-2 is 48:
	 * pod 2: containerA, QPS = 50, Response time = 0	weighted = 0
	 * pod 2: containerB, QPS = 1, Response time = 288	weighted = 288
	 * pod 3, containerC, QPS = 80, Response time = 75	weighted = 6000
	 *
	 * Total transactions = 131, total response time = 6288, average = 48
	 */

	 for _, dto := range dtoList {
	 	if ("service-2" == *dto.Id) {
	 		for _, comm := range dto.CommoditiesSold {
	 			if (*comm.CommodityType == proto.CommodityDTO_RESPONSE_TIME) {
	 					assert.Equal(t, *comm.Used, 48.0)
	 					return
				}
			}
		}
	}
	 t.Error("Could not locate ResponseTime commodity sold in service-2")
}
