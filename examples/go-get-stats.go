package bigip

import (
	"fmt"
	"github.com/kzittritsch/go-bigip"
)

func main() {
	// Connect to the BIG-IP system.
	f5 := bigip.NewSession("ltm.example.com", "yadda", "yadda")

	// Get all pools
	pools, err := f5.Pools()
	if err != nil {
		fmt.Println(err)
	}
	// get pool name and active members for all pools
	for _, pool := range pools.Pools {
		PoolStats, err := f5.GetPoolStatistics(pool.FullPath)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Pool name: ", PoolStats.Entries.Tmname)
		fmt.Println("Pool Active Members: ", PoolStats.Entries.Activemembercnt)
	}

	// get all virtual servers
	virtuals, err := f5.VirtualServers()
	if err != nil {
		fmt.Println(err)
	}
	// get name, destination, enabled state for all vss
	for _, vss := range virtuals.VirtualServers {
		VSSStats, err := f5.GetVSSStatistics(vss.FullPath)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("VSS Name: ", VSSStats.Entries.Tmname)
		fmt.Println("VSS Destination: ", VSSStats.Entries.Destination)
		fmt.Println("VSS Status enabled state: ", VSSStats.Entries.StatusEnabledstate)
	}
}

