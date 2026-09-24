package query_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/query"
)

var _ networks.ListOptsBuilder = (*query.ListOpts)(nil)
var _ ports.ListOptsBuilder = (*query.ListOpts)(nil)
var _ subnets.ListOptsBuilder = (*query.ListOpts)(nil)

func ExampleListOpts_And_by_id() {
	q := query.New(ports.ListOpts{
		Name: "Jules",
	}).And("id", "123", "321", "12345")
	fmt.Println(q)
	//Output: ?id=123&id=321&id=12345&name=Jules
}

func ExampleListOpts_And_by_name() {
	q := query.New(ports.ListOpts{}).
		And("name", "port-1", "port-&321", "the-other-port")
	fmt.Println(q)
	//Output: ?name=port-1&name=port-%26321&name=the-other-port
}

func ExampleListOpts_And_by_Name_and_tag() {
	q := query.New(ports.ListOpts{}).
		And("name", "port-1", "port-3").
		And("tags", "my-tag")
	fmt.Println(q)
	//Output: ?name=port-1&name=port-3&tags=my-tag
}

func ExampleListOpts_And_by_id_twice() {
	q := query.New(ports.ListOpts{}).
		And("id", "1", "2", "3").
		And("id", "2", "3", "4")
	fmt.Println(q)
	//Output: ?id=2&id=3
}

func ExampleListOpts_And_by_id_twice_plus_ListOpts() {
	q := query.New(ports.ListOpts{ID: "3"}).
		And("id", "1", "2", "3").
		And("id", "3", "4", "5")
	fmt.Println(q)
	//Output: ?id=3
}

func TestToPortListQuery(t *testing.T) {
	for _, tc := range [...]struct {
		name          string
		base          interface{}
		andProperty   string
		andItems      []interface{}
		expected      string
		expectedError bool
	}{
		{
			"valid",
			ports.ListOpts{},
			"name",
			[]interface{}{"port-1"},
			"?name=port-1",
			false,
		},
		{
			"invalid_field",
			ports.ListOpts{},
			"door",
			[]interface{}{"pod bay"},
			"?",
			true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q, err := query.New(tc.base).And(tc.andProperty, tc.andItems...).ToPortListQuery()
			if q != tc.expected {
				t.Errorf("expected query %q, got %q", tc.expected, q)
			}
			if (err != nil) != tc.expectedError {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else {
					t.Errorf("expected error, got nil")
				}
			}
		})
	}
}
