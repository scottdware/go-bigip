package bigip

import (
	"fmt"

	"github.com/scottdware/go-bigip"
)

// Note: Had to call main something different as the package complained main was in LTM's example
func mainGTM() {
	// Connect to the BIG-IP system.
	f5 := bigip.NewSession("gtm.company.com", "admin", "secret", nil)

	// Get a list of all WideIP's, and print their names to the console.
	result, err := f5.GetGTMWideIPs(bigip.ARecord)
	if err != nil {
		fmt.Println(err)
	}

	for _, wideip := range result.GTMWideIPs {
		fmt.Println(wideip.Name)
	}

	// Ordering here just makes things easier.
	// Since Pools can exist without any other reliance, we start there.
	poolConfig := &bigip.GTMAPool{
		Name:      "sample.company.com_pool",
		Partition: "Common",
	}
	f5.AddGTMAPool(poolConfig)

	// Create a WideIP Pool config using our Pool Config
	wipPoolConfig := bigip.GTMWideIPPool{
		Name:      poolConfig.Name,
		Partition: poolConfig.Partition,
	}

	// right now we only have a 1:1 mapping - but WideIPs can have multiple pools
	wipPools := &[]bigip.GTMWideIPPool{wipPoolConfig}

	wipConfig := &bigip.GTMWideIP{
		Name:      "sample.company.com",
		Partition: "Common",
		Pools:     wipPools,
	}

	// Crate the new WideIP
	f5.AddGTMWideIP(wipConfig, bigip.ARecord)

	// Add Pools
	// This part is faked out - and you will have to provide some information to make this work.
	// Namely you have to know the LTM Server which was setup in the GTM (with VirtualServerDiscovery turned on)
	// And you will also need at least one virtual server by full path  e.g.  /partition/name

	fullPathAPool := fmt.Sprintf("/%s/%s", poolConfig.Partition, poolConfig.Name)
	virtualServerPath := "/Common/baseapp_80_vs"
	serverPath := "/Common/someltm"

	// Add the Virtual Server to the list of members for the pool
	f5.CreateGTMAPoolMember(fullPathAPool, serverPath, virtualServerPath)

	// Unwind all the work
	f5.DeleteGTMAPoolMember(fullPathAPool, serverPath, virtualServerPath)

	wipFullPath := fmt.Sprintf("/%s/%s", wipConfig.Partition, wipConfig.Name)
	f5.DeleteGTMWideIP(wipFullPath, bigip.ARecord)

	f5.DeleteGTMPool(fullPathAPool, bigip.ARecord)
}
