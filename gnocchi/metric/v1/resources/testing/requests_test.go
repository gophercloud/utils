package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resources"
	fake "github.com/gophercloud/utils/gnocchi/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/resource/generic", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ResourceListResult)
	})

	count := 0

	resources.List(fake.ServiceClient(), resources.ListOpts{}, "").EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := resources.ExtractResources(page)
		if err != nil {
			t.Errorf("Failed to extract resources: %v", err)
			return false, nil
		}

		expected := []resources.Resource{
			Resource1,
			Resource2,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/resource/compute_instance_network/75274f99-faf6-4112-a6d5-2794cb07c789", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ResourceGetResult)
	})

	s, err := resources.Get(fake.ServiceClient(), "75274f99-faf6-4112-a6d5-2794cb07c789", "compute_instance_network").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.CreatedByProjectID, "3d40ca37723449118987b9f288f4ae84")
	th.AssertEquals(t, s.CreatedByUserID, "fdcfb420c09645e69e177a0bb1950884")
	th.AssertEquals(t, s.Creator, "fdcfb420c09645e69e177a0bb1950884:3d40ca37723449118987b9f288f4ae84")
	th.AssertEquals(t, s.ID, "75274f99-faf6-4112-a6d5-2794cb07c789")
	th.AssertDeepEquals(t, s.Metrics, map[string]string{
		"network.incoming.bytes.rate":   "01b2953e-de74-448a-a305-c84440697933",
		"network.outgoing.bytes.rate":   "4ac0041b-3bf7-441d-a95a-d3e2f1691158",
		"network.incoming.packets.rate": "5a64328e-8a7c-4c6a-99df-2e6d17440142",
		"network.outgoing.packets.rate": "dc9f3198-155b-4b88-a92c-58a3853ce2b2",
	})
	th.AssertEquals(t, s.OriginalResourceID, "75274f99-faf6-4112-a6d5-2794cb07c789")
	th.AssertEquals(t, s.ProjectID, "4154f08883334e0494c41155c33c0fc9")
	th.AssertEquals(t, s.RevisionStart, "2018-01-01T11:44:31.742031+00:00")
	th.AssertEquals(t, s.RevisionEnd, "")
	th.AssertEquals(t, s.StartedAt, "2018-01-01T11:44:31.742011+00:00")
	th.AssertEquals(t, s.EndedAt, "")
	th.AssertEquals(t, s.Type, "compute_instance_network")
	th.AssertEquals(t, s.UserID, "bd5874d666624b24a9f01c128871e4ac")
}
