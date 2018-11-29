package api

import (
	"fmt"
	"strings"
)

type WorkerNodePools []WorkerNodePool

// NodePoolAvailabilityZoneDependencies produces a list WorkerNodePools which need to be rolled before the named
// pool when using the 'AvailabilityZone' NodePoolRollingStrategy.
// Only other pools using the 'AvailabilityZone' strategy are considered.
// Only nodepools containing a single AZ can use this strategy!
// Returns a comma separated quoted list, e.g. "pool1","pool2","pool3"
func (pools WorkerNodePools) NodePoolAvailabilityZoneDependencies(nodePool WorkerNodePool, subnets Subnets) (string, error) {
	var order []string
	order, err := pools.rolloutAZOrder()
	if err != nil {
		return "", fmt.Errorf("can't resolve nodepool availability zone ordering dependencies: %v", err)
	}

	position, err := azPosition(order, nodePool.Subnets[0].AvailabilityZone)
	if err != nil {
		return "", err
	}
	if position == 0 {
		// a nodePool with position 0 doesn't have any other nodepool dependencies
		return "", nil
	}

	return `"` + strings.Join(pools.allNodePoolsinAZ(order[position-1]), `","`) + `"`, nil
}

// createAZOrder works out an order for availability zones from the order in which nodepools and subnets have been ordered in the cluster.yaml.
// It also provides some validation - preventing users from selecting subnets in separate az's which would be impossible to order by placing dependencies
// at the nodepool stack level.
func (pools WorkerNodePools) rolloutAZOrder() ([]string, error) {
	var azOrder []string
	seen := make(map[string]bool)
	for _, pool := range pools {
		if pool.NodePoolRollingStrategy == "AvailabilityZone" {
			if len(pool.Subnets) == 0 {
				return azOrder, fmt.Errorf("worker nodepool %s has 'AvailabilityZone' rolling strategy but has no subnets", pool.NodePoolName)
			}
			azCheck := pool.Subnets[0].AvailabilityZone
			for _, subnet := range pool.Subnets {
				if subnet.AvailabilityZone == "" {
					return azOrder, fmt.Errorf("worker nodepool %s can not use the 'AvailabilityZone' rolling strategy because its subnet %s has an empty availability zone", pool.NodePoolName, subnet.Name)
				}
				if subnet.AvailabilityZone != azCheck {
					return azOrder, fmt.Errorf("worker nodepool %s can't have subnets in different availability zones and also use the 'AvailabilityZone' rolling strategy", pool.NodePoolName)
				}
				if !seen[subnet.AvailabilityZone] {
					azOrder = append(azOrder, subnet.AvailabilityZone)
					seen[subnet.AvailabilityZone] = true
				}
			}
		}
	}
	return azOrder, nil
}

// azPosition tells us which integer position a particular az lies within an list
func azPosition(azList []string, az string) (int, error) {
	for index, value := range azList {
		if value == az {
			return index, nil
		}
	}
	return 0, fmt.Errorf("could not find az %s in the azorder list: %v", az, azList)
}

// allNodePoolsinAZ returns a slice of nodepool names with subnets within a given availability zone
func (pools WorkerNodePools) allNodePoolsinAZ(az string) []string {
	var poolNames []string
	for _, pool := range pools {
		if pool.Subnets[0].AvailabilityZone == az {
			poolNames = append(poolNames, pool.NodePoolName)
		}
	}
	return poolNames
}
