package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/utils/gnocchi/metric/v1/common"
	"github.com/gophercloud/utils/gnocchi/metric/v1/metrics"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/metric", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, MetricsListResult)
	})

	count := 0

	metrics.List(fake.ServiceClient(), metrics.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := metrics.ExtractMetrics(page)
		if err != nil {
			t.Errorf("Failed to extract metrics: %v", err)
			return false, nil
		}

		expected := []metrics.Metric{
			Metric1,
			Metric2,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
