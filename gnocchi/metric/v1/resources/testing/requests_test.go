package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/utils/gnocchi/metric/v1/common"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resources"
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
