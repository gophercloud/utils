package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/utils/gnocchi/metric/v1/measures"
	fake "github.com/gophercloud/utils/gnocchi/testhelper/client"
)

func TestListMeasures(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/metric/9e5a6441-1044-4181-b66e-34e180753040/measures", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, MeasuresListResult)
	})

	metricID := "9e5a6441-1044-4181-b66e-34e180753040"
	startTime := time.Date(2018, 1, 10, 12, 0, 0, 0, time.UTC)
	stopTime := time.Date(2018, 1, 10, 14, 0, 5, 0, time.UTC)
	opts := measures.ListOpts{
		Start:       &startTime,
		Stop:        &stopTime,
		Granularity: "1h",
	}
	expected := ListMeasuresExpected
	pages := 0
	err := measures.List(fake.ServiceClient(), metricID, opts).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := measures.ExtractMeasures(page)
		th.AssertNoErr(t, err)

		if len(actual) != 3 {
			t.Fatalf("Expected 2 measures, got %d", len(actual))
		}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, pages)
}

func TestCreateMeasures(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/metric/9e5a6441-1044-4181-b66e-34e180753040/measures", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json, */*")
		w.WriteHeader(http.StatusAccepted)
	})

	firstMeasureTimestamp := time.Date(2018, 1, 18, 12, 31, 0, 0, time.UTC)
	secondMeasureTimestamp := time.Date(2018, 1, 18, 14, 32, 0, 0, time.UTC)
	createOpts := measures.CreateOpts{
		Measures: []measures.MeasureOpts{
			{
				Timestamp: &firstMeasureTimestamp,
				Value:     101.2,
			},
			{
				Timestamp: &secondMeasureTimestamp,
				Value:     102,
			},
		},
	}
	res := measures.Create(fake.ServiceClient(), "9e5a6441-1044-4181-b66e-34e180753040", createOpts)
	if res.Err.Error() == "EOF" {
		res.Err = nil
	}
	th.AssertNoErr(t, res.Err)
}
