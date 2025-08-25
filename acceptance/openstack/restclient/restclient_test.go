// +build acceptance restclient

package restclient

import (
	"testing"

	acc_clients "github.com/gophercloud/gophercloud/acceptance/clients"
	acc_tools "github.com/gophercloud/gophercloud/acceptance/tools"

	th "github.com/gophercloud/gophercloud/testhelper"
	cc "github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/gophercloud/utils/openstack/restclient"
)

func TestRESTClient(t *testing.T) {
	acc_clients.RequireAdmin(t)

	// This will be populated by environment variables.
	clientOpts := &cc.ClientOpts{}

	computeClient, err := cc.NewServiceClient("compute", clientOpts)
	th.AssertNoErr(t, err)

	// Test creating a flavor
	flavorName := acc_tools.RandomString("TESTACC-", 8)
	flavorID := acc_tools.RandomString("TESTACC-", 8)
	flavorOpts := map[string]interface{}{
		"name":  flavorName,
		"ram":   512,
		"vcpus": 1,
		"disk":  5,
		"id":    flavorID,
	}

	postOpts := &restclient.PostOpts{
		Params: map[string]interface{}{"flavor": flavorOpts},
	}

	postURL := computeClient.ServiceURL("flavors")
	postRes := restclient.Post(computeClient, postURL, postOpts)
	th.AssertNoErr(t, postRes.Err)
	flavorResult, err := postRes.Extract()
	th.AssertNoErr(t, postRes.Err)
	acc_tools.PrintResource(t, flavorResult)

	// Test deleting a flavor
	defer func() {
		deleteURL := computeClient.ServiceURL("flavors", flavorID)
		deleteRes := restclient.Delete(computeClient, deleteURL, nil)
		th.AssertNoErr(t, deleteRes.Err)
		err = deleteRes.ExtractErr()
		th.AssertNoErr(t, err)
	}()

	// Test retrieving a flavor
	getURL := computeClient.ServiceURL("flavors", flavorID)
	getRes := restclient.Get(computeClient, getURL, nil)
	th.AssertNoErr(t, getRes.Err)

	flavorResult, err = getRes.Extract()
	th.AssertNoErr(t, err)

	flavor := flavorResult["flavor"].(map[string]interface{})

	acc_tools.PrintResource(t, flavor)

	th.AssertEquals(t, flavor["disk"], float64(5))
	th.AssertEquals(t, flavor["id"], flavorID)
	th.AssertEquals(t, flavor["name"], flavorName)
	th.AssertEquals(t, flavor["ram"], float64(512))
	th.AssertEquals(t, flavor["swap"], "")
	th.AssertEquals(t, flavor["vcpus"], float64(1))

	// Test listing flavors
	getOpts := &restclient.GetOpts{
		Query: map[string]interface{}{
			"limit": 2,
		},
	}

	getURL = computeClient.ServiceURL("flavors")
	getRes = restclient.Get(computeClient, getURL, getOpts)
	th.AssertNoErr(t, getRes.Err)
	flavorResult, err = getRes.Extract()
	th.AssertNoErr(t, err)

	flavors := flavorResult["flavors"].([]interface{})
	acc_tools.PrintResource(t, flavors)
	th.AssertEquals(t, len(flavors), 2)
}
