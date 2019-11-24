// +build acceptance clientconfig

package clientconfig

import (
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	acc_compute "github.com/gophercloud/gophercloud/acceptance/openstack/compute/v2"
	acc_tools "github.com/gophercloud/gophercloud/acceptance/tools"

	cc "github.com/gophercloud/utils/openstack/clientconfig"
)

func TestServerCreateDestroy(t *testing.T) {
	// This will be populated by environment variables.
	clientOpts := &cc.ClientOpts{}

	client, err := cc.NewServiceClient("compute", clientOpts)
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	server, err := acc_compute.CreateServer(t, client)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	defer acc_compute.DeleteServer(t, client, server)

	newServer, err := servers.Get(client, server.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get server %s: %v", server.ID, err)
	}

	acc_tools.PrintResource(t, newServer)
}

func TestEndpointType(t *testing.T) {
	clientOpts := &cc.ClientOpts{
		EndpointType: "admin",
	}
	client, err := cc.NewServiceClient("identity", clientOpts)
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	if !strings.Contains(client.Endpoint, "35357") {
		t.Fatalf("Endpoint was not correctly set to admin interface")
	}
}
