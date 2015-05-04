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

	// Create a couple of nodes.
	b.CreateNode("web-server-1", "192.168.1.50")
	b.CreateNode("web-server-2", "192.168.1.51")
	b.CreateNode("ssl-web-server-1", "10.2.2.50")
	b.CreateNode("ssl-web-server-2", "10.2.2.51")

	// Create a pool, and add members to it. When adding a member, you must
	// specify the port in the format of <node name>:<port>.
	b.CreatePool("web_farm_80_pool")
	b.AddPoolMember("web_farm_80_pool", "web-server-1:80")
	b.AddPoolMember("web_farm_80_pool", "web-server-2:80")

	b.CreatePool("ssl_443_pool")
	b.AddPoolMember("ssl_443_pool", "ssl-web-server-1:443")
	b.AddPoolMember("ssl_443_pool", "ssl-web-server-2:443")

	// Create a virtual server, with the above pool. The third field is the subnet
	// mask, and that can either be in CIDR notation or decimal. For any/all destinations
	// and ports, use '0' for the mask and/or port.
	b.CreateVirtualServer("web_farm_VS", "0.0.0.0", "0.0.0.0", "web_farm_80_pool", 80)
	b.CreateVirtualServer("ssl_web_farm_VS", "10.1.1.0", "24", "ssl_443_pool", 443)

	// Remove a pool member.
	b.DeletePoolMember("web_farm_80_pool", "web-server-2:80")
}
