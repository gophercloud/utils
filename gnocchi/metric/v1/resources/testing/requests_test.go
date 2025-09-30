package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/utils/v2/gnocchi/metric/v1/resources"
	fake "github.com/gophercloud/utils/v2/gnocchi/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/generic", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := r.ParseForm()
		th.AssertNoErr(t, err)
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprint(w, ResourceListResult)
		case "789a7f65-977d-40f4-beed-f717100125f5":
			fmt.Fprint(w, `[]`)
		default:
			t.Fatalf("/v1/resources invoked with unexpected marker=[%s]", marker)
		}
	})

	count := 0

	err := resources.List(fake.ServiceClient(fakeServer), resources.ListOpts{}, "generic").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/compute_instance_network/75274f99-faf6-4112-a6d5-2794cb07c789", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ResourceGetResult)
	})

	s, err := resources.Get(context.TODO(), fake.ServiceClient(fakeServer), "compute_instance_network", "75274f99-faf6-4112-a6d5-2794cb07c789").Extract()
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
	th.AssertEquals(t, s.RevisionStart, time.Date(2018, 01, 01, 11, 44, 31, 742031000, time.UTC))
	th.AssertEquals(t, s.RevisionEnd, time.Time{})
	th.AssertEquals(t, s.StartedAt, time.Date(2018, 01, 01, 11, 44, 31, 742011000, time.UTC))
	th.AssertEquals(t, s.EndedAt, time.Time{})
	th.AssertEquals(t, s.Type, "compute_instance_network")
	th.AssertEquals(t, s.UserID, "bd5874d666624b24a9f01c128871e4ac")
	th.AssertDeepEquals(t, s.ExtraAttributes, map[string]any{
		"iface_name": "eth0",
	})
}

func TestCreateWithoutMetrics(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/generic", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceCreateWithoutMetricsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, ResourceCreateWithoutMetricsResult)
	})

	opts := resources.CreateOpts{
		ID:        "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55",
		ProjectID: "4154f088-8333-4e04-94c4-1155c33c0fc9",
		UserID:    "bd5874d6-6662-4b24-a9f01c128871e4ac",
	}
	s, err := resources.Create(context.TODO(), fake.ServiceClient(fakeServer), "generic", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.CreatedByProjectID, "3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.CreatedByUserID, "fdcfb420-c096-45e6-9e177a0bb1950884")
	th.AssertEquals(t, s.Creator, "fdcfb420-c096-45e6-9e177a0bb1950884:3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.ID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertDeepEquals(t, s.Metrics, map[string]string{})
	th.AssertEquals(t, s.OriginalResourceID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertEquals(t, s.ProjectID, "4154f088-8333-4e04-94c4-1155c33c0fc9")
	th.AssertEquals(t, s.RevisionStart, time.Date(2018, 1, 3, 11, 44, 31, 155773000, time.UTC))
	th.AssertEquals(t, s.RevisionEnd, time.Time{})
	th.AssertEquals(t, s.StartedAt, time.Date(2018, 1, 3, 11, 44, 31, 155732000, time.UTC))
	th.AssertEquals(t, s.EndedAt, time.Time{})
	th.AssertEquals(t, s.Type, "generic")
	th.AssertEquals(t, s.UserID, "bd5874d6-6662-4b24-a9f01c128871e4ac")
}

func TestCreateLinkMetrics(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/compute_instance_network", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceCreateLinkMetricsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, ResourceCreateLinkMetricsResult)
	})

	startedAt := time.Date(2018, 1, 2, 23, 23, 34, 0, time.UTC)
	endedAt := time.Date(2018, 1, 4, 10, 00, 12, 0, time.UTC)
	opts := resources.CreateOpts{
		ID:        "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55",
		ProjectID: "4154f088-8333-4e04-94c4-1155c33c0fc9",
		UserID:    "bd5874d6-6662-4b24-a9f01c128871e4ac",
		StartedAt: &startedAt,
		EndedAt:   &endedAt,
		Metrics: map[string]any{
			"network.incoming.bytes.rate": "01b2953e-de74-448a-a305-c84440697933",
			"network.outgoing.bytes.rate": "dc9f3198-155b-4b88-a92c-58a3853ce2b2",
		},
	}
	s, err := resources.Create(context.TODO(), fake.ServiceClient(fakeServer), "compute_instance_network", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.CreatedByProjectID, "3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.CreatedByUserID, "fdcfb420-c096-45e6-9e177a0bb1950884")
	th.AssertEquals(t, s.Creator, "fdcfb420-c096-45e6-9e177a0bb1950884:3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.ID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertDeepEquals(t, s.Metrics, map[string]string{
		"network.incoming.bytes.rate": "01b2953e-de74-448a-a305-c84440697933",
		"network.outgoing.bytes.rate": "dc9f3198-155b-4b88-a92c-58a3853ce2b2",
	})
	th.AssertEquals(t, s.OriginalResourceID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertEquals(t, s.ProjectID, "4154f088-8333-4e04-94c4-1155c33c0fc9")
	th.AssertEquals(t, s.RevisionStart, time.Date(2018, 1, 2, 23, 23, 34, 155813000, time.UTC))
	th.AssertEquals(t, s.RevisionEnd, time.Time{})
	th.AssertEquals(t, s.StartedAt, time.Date(2018, 1, 2, 23, 23, 34, 0, time.UTC))
	th.AssertEquals(t, s.EndedAt, time.Date(2018, 1, 4, 10, 00, 12, 0, time.UTC))
	th.AssertEquals(t, s.Type, "compute_instance_network")
	th.AssertEquals(t, s.UserID, "bd5874d6-6662-4b24-a9f01c128871e4ac")
}

func TestCreateWithMetrics(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/compute_instance_disk", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceCreateWithMetricsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, ResourceCreateWithMetricsResult)
	})

	endedAt := time.Date(2018, 1, 9, 20, 0, 0, 0, time.UTC)
	opts := resources.CreateOpts{
		ID:        "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55",
		ProjectID: "4154f088-8333-4e04-94c4-1155c33c0fc9",
		UserID:    "bd5874d6-6662-4b24-a9f01c128871e4ac",
		EndedAt:   &endedAt,
		Metrics: map[string]any{
			"disk.write.bytes.rate": map[string]string{
				"archive_policy_name": "high",
			},
		},
	}
	s, err := resources.Create(context.TODO(), fake.ServiceClient(fakeServer), "compute_instance_disk", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.CreatedByProjectID, "3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.CreatedByUserID, "fdcfb420-c096-45e6-9e177a0bb1950884")
	th.AssertEquals(t, s.Creator, "fdcfb420-c096-45e6-9e177a0bb1950884:3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.ID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertDeepEquals(t, s.Metrics, map[string]string{
		"disk.write.bytes.rate": "0a2da84d-4753-43f5-a65f-0f8d44d2766c",
	})
	th.AssertEquals(t, s.OriginalResourceID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertEquals(t, s.ProjectID, "4154f088-8333-4e04-94c4-1155c33c0fc9")
	th.AssertEquals(t, s.RevisionStart, time.Date(2018, 1, 2, 23, 23, 34, 155813000, time.UTC))
	th.AssertEquals(t, s.RevisionEnd, time.Time{})
	th.AssertEquals(t, s.StartedAt, time.Date(2018, 1, 2, 23, 23, 34, 155773000, time.UTC))
	th.AssertEquals(t, s.EndedAt, time.Date(2018, 1, 9, 20, 00, 00, 0, time.UTC))
	th.AssertEquals(t, s.Type, "compute_instance_disk")
	th.AssertEquals(t, s.UserID, "bd5874d6-6662-4b24-a9f01c128871e4ac")
}

func TestUpdateLinkMetrics(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/compute_instance_network/23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceUpdateLinkMetricsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ResourceUpdateLinkMetricsResponse)
	})

	endedAt := time.Date(2018, 1, 14, 13, 0, 0, 0, time.UTC)
	metrics := map[string]any{
		"network.incoming.bytes.rate": "01b2953e-de74-448a-a305-c84440697933",
	}
	updateOpts := resources.UpdateOpts{
		EndedAt: &endedAt,
		Metrics: &metrics,
	}
	s, err := resources.Update(context.TODO(), fake.ServiceClient(fakeServer), "compute_instance_network", "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55", updateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.CreatedByProjectID, "3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.CreatedByUserID, "fdcfb420-c096-45e6-9e177a0bb1950884")
	th.AssertEquals(t, s.Creator, "fdcfb420-c096-45e6-9e177a0bb1950884:3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.ID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertDeepEquals(t, s.Metrics, map[string]string{
		"network.incoming.bytes.rate": "01b2953e-de74-448a-a305-c84440697933",
	})
	th.AssertEquals(t, s.OriginalResourceID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertEquals(t, s.ProjectID, "4154f088-8333-4e04-94c4-1155c33c0fc9")
	th.AssertEquals(t, s.RevisionStart, time.Date(2018, 1, 12, 13, 44, 34, 742031000, time.UTC))
	th.AssertEquals(t, s.RevisionEnd, time.Time{})
	th.AssertEquals(t, s.StartedAt, time.Date(2018, 1, 12, 13, 44, 34, 742011000, time.UTC))
	th.AssertEquals(t, s.EndedAt, time.Date(2018, 1, 14, 13, 0, 0, 0, time.UTC))
	th.AssertEquals(t, s.Type, "compute_instance_network")
	th.AssertEquals(t, s.UserID, "bd5874d6-6662-4b24-a9f01c128871e4ac")
}

func TestUpdateCreateMetrics(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/compute_instance_network/23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceUpdateCreateMetricsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ResourceUpdateCreateMetricsResponse)
	})

	startedAt := time.Date(2018, 1, 12, 11, 0, 0, 0, time.UTC)
	metrics := map[string]any{
		"disk.read.bytes.rate": map[string]string{
			"archive_policy_name": "low",
		},
	}
	updateOpts := resources.UpdateOpts{
		StartedAt: &startedAt,
		Metrics:   &metrics,
	}
	s, err := resources.Update(context.TODO(), fake.ServiceClient(fakeServer), "compute_instance_network", "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55", updateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.CreatedByProjectID, "3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.CreatedByUserID, "fdcfb420-c096-45e6-9e177a0bb1950884")
	th.AssertEquals(t, s.Creator, "fdcfb420-c096-45e6-9e177a0bb1950884:3d40ca37-7234-4911-8987b9f288f4ae84")
	th.AssertEquals(t, s.ID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertDeepEquals(t, s.Metrics, map[string]string{
		"disk.read.bytes.rate": "ed1bb76f-6ccc-4ad2-994c-dbb19ddccbae",
	})
	th.AssertEquals(t, s.OriginalResourceID, "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertEquals(t, s.ProjectID, "4154f088-8333-4e04-94c4-1155c33c0fc9")
	th.AssertEquals(t, s.RevisionStart, time.Date(2018, 1, 12, 12, 00, 34, 742031000, time.UTC))
	th.AssertEquals(t, s.RevisionEnd, time.Time{})
	th.AssertEquals(t, s.StartedAt, time.Date(2018, 1, 12, 11, 00, 00, 0, time.UTC))
	th.AssertEquals(t, s.EndedAt, time.Time{})
	th.AssertEquals(t, s.Type, "compute_instance_disk")
	th.AssertEquals(t, s.UserID, "bd5874d6-6662-4b24-a9f01c128871e4ac")
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource/generic/23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := resources.Delete(context.TODO(), fake.ServiceClient(fakeServer), "generic", "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55")
	th.AssertNoErr(t, res.Err)
}
