package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/utils/v2/gnocchi/metric/v1/status"
	fake "github.com/gophercloud/utils/v2/gnocchi/testhelper/client"
)

func TestGetWithDetails(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, StatusGetWithDetailsResult)
	})

	details := true

	getOpts := status.GetOpts{
		Details: &details,
	}

	s, err := status.Get(context.TODO(), fake.ServiceClient(fakeServer), getOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s.Metricd, GetStatusWithDetailsExpected.Metricd)
	th.AssertDeepEquals(t, s.Storage, GetStatusWithDetailsExpected.Storage)
}

func TestGetWithoutDetails(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, StatusGetWithoutDetailsResult)
	})

	details := false

	getOpts := status.GetOpts{
		Details: &details,
	}

	s, err := status.Get(context.TODO(), fake.ServiceClient(fakeServer), getOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s.Metricd, GetStatusWithoutDetailsExpected.Metricd)
	th.AssertDeepEquals(t, s.Storage, GetStatusWithoutDetailsExpected.Storage)
}
