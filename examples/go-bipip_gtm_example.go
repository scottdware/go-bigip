package bigip

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	bigip "github.com/wwt/go-bigip"
)

// StartTest just tests our f5/gtm stuffs
func StartTest() {

	host := os.Getenv("GTM_HOST")
	user := os.Getenv("GTM_USER")
	pass := os.Getenv("GTM_PASS")

	options := bigip.ConfigOptions{}
	// we use AD integration for auth so we need to use tmos
	// else use NewSession(host, user, pass, &options)
	f5, err := bigip.NewTokenSession(host, user, pass, "tmos", &options)
	if err != nil {
		fmt.Printf("Nope No Workie\n%+v\n", err)
	}

	things, err := f5.GetWideIPs(bigip.ARecord)
	if err != nil {
		fmt.Printf("Nope No Workie\n%+v\n", err)
	}

	spew.Dump(things)

}
