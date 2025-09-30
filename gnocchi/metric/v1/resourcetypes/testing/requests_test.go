package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/utils/v2/gnocchi/metric/v1/resourcetypes"
	fake "github.com/gophercloud/utils/v2/gnocchi/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource_type", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ResourceTypeListResult)
	})

	count := 0

	err := resourcetypes.List(fake.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := resourcetypes.ExtractResourceTypes(page)
		if err != nil {
			t.Errorf("Failed to extract resource types: %v", err)
			return false, nil
		}

		expected := []resourcetypes.ResourceType{
			ResourceType1,
			ResourceType2,
			ResourceType3,
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

	fakeServer.Mux.HandleFunc("/v1/resource_type/compute_instance", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ResourceTypeGetResult)
	})

	s, err := resourcetypes.Get(context.TODO(), fake.ServiceClient(fakeServer), "compute_instance").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "compute_instance")
	th.AssertEquals(t, s.State, "active")
	th.AssertDeepEquals(t, s.Attributes, map[string]resourcetypes.Attribute{
		"host": {
			Type: "string",
			Details: map[string]any{
				"max_length": float64(255),
				"min_length": float64(0),
				"required":   true,
			},
		},
		"image_ref": {
			Type: "uuid",
			Details: map[string]any{
				"required": false,
			},
		},
	})
}

func TestCreateWithoutAttributes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource_type", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceTypeCreateWithoutAttributesRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, ResourceTypeCreateWithoutAttributesResult)
	})

	opts := resourcetypes.CreateOpts{
		Name: "identity_project",
	}
	s, err := resourcetypes.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "identity_project")
	th.AssertEquals(t, s.State, "active")
	th.AssertDeepEquals(t, s.Attributes, map[string]resourcetypes.Attribute{})
}

func TestCreateWithAttributes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource_type", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceTypeCreateWithAttributesRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, ResourceTypeCreateWithAttributesResult)
	})

	opts := resourcetypes.CreateOpts{
		Name: "compute_instance_network",
		Attributes: map[string]resourcetypes.AttributeOpts{
			"port_name": {
				Type: "string",
				Details: map[string]any{
					"max_length": 128,
					"required":   false,
				},
			},
			"port_id": {
				Type: "uuid",
				Details: map[string]any{
					"required": true,
				},
			},
		},
	}
	s, err := resourcetypes.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "compute_instance_network")
	th.AssertEquals(t, s.State, "active")
	th.AssertDeepEquals(t, s.Attributes, map[string]resourcetypes.Attribute{
		"port_name": {
			Type: "string",
			Details: map[string]any{
				"max_length": float64(128),
				"min_length": float64(0),
				"required":   false,
			},
		},
		"port_id": {
			Type: "uuid",
			Details: map[string]any{
				"required": true,
			},
		},
	})
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource_type/identity_project", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json-patch+json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceTypeUpdateRequest)

		w.Header().Add("Content-Type", "application/json-patch+json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ResourceTypeUpdateResult)
	})

	enabledAttributeOptions := resourcetypes.AttributeOpts{
		Details: map[string]any{
			"required": true,
			"options": map[string]any{
				"fill": true,
			},
		},
		Type: "bool",
	}
	parendIDAttributeOptions := resourcetypes.AttributeOpts{
		Details: map[string]any{
			"required": false,
		},
		Type: "uuid",
	}
	opts := resourcetypes.UpdateOpts{
		Attributes: []resourcetypes.AttributeUpdateOpts{
			{
				Name:      "enabled",
				Operation: resourcetypes.AttributeAdd,
				Value:     &enabledAttributeOptions,
			},
			{
				Name:      "parent_id",
				Operation: resourcetypes.AttributeAdd,
				Value:     &parendIDAttributeOptions,
			},
			{
				Name:      "domain_id",
				Operation: resourcetypes.AttributeRemove,
			},
		},
	}

	s, err := resourcetypes.Update(context.TODO(), fake.ServiceClient(fakeServer), "identity_project", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "identity_project")
	th.AssertEquals(t, s.State, "active")
	th.AssertDeepEquals(t, s.Attributes, map[string]resourcetypes.Attribute{
		"enabled": {
			Type: "bool",
			Details: map[string]any{
				"required": true,
			},
		},
		"parent_id": {
			Type: "uuid",
			Details: map[string]any{
				"required": false,
			},
		},
		"name": {
			Type: "string",
			Details: map[string]any{
				"required":   true,
				"min_length": float64(0),
				"max_length": float64(128),
			},
		},
	})
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/resource_type/compute_instance_network", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := resourcetypes.Delete(context.TODO(), fake.ServiceClient(fakeServer), "compute_instance_network")
	th.AssertNoErr(t, res.Err)
}
