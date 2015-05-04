package bigip

import (
	"fmt"
	"github.com/scottdware/go-bigip"
)

func main() {
	// Connect to the BIG-IP system.
	b := bigip.NewSession("ltm.company.com", "admin", "secret")

	// Get a list of all VLAN's, and print their names to the console.
	vlans, err := b.Vlans()
	if err != nil {
		fmt.Println(err)
	}

	for _, vlan := range vlans.Vlans {
		fmt.Println(vlan.Name)
	}

	// Create a VLAN
	b.CreateVlan("vlan1138")

	// Add an untagged interface to a VLAN.
	b.AddInterfaceToVlan("vlan1138", "1.2", false)

	// Delete a VLAN.
	b.DeleteVlan("vlan1138")
}
