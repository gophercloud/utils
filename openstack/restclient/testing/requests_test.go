package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gophercloud/gophercloud"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/utils/openstack/restclient"
)

func TestBuildQueryString(t *testing.T) {
	testCases := map[string]interface{}{
		"a": 2,
		"b": "foo",
		"c": true,
		"d": []string{"one", "two", "three"},
		"e": []int{1, 2, 3},
		"f": map[string]string{"foo": "bar"},
		"g": false,
	}

	expected := &url.URL{RawQuery: "a=2&b=foo&c=true&d=one&d=two&d=three&e=1&e=2&e=3&f=%7B%27foo%27%3A%27bar%27%7D&g=false"}

	actual, err := restclient.BuildQueryString(testCases)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expected, actual)
}

func TestBasic(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"foo": "bar"}`)
	})

	url := fmt.Sprintf("%s/route", th.Endpoint())

	c := new(gophercloud.ServiceClient)
	c.ProviderClient = new(gophercloud.ProviderClient)

	expected := map[string]interface{}{
		"foo": "bar",
	}

	// shared params
	params := map[string]interface{}{
		"bar": "baz",
	}

	// Get
	getRes := restclient.Get(c, url, nil)
	th.AssertNoErr(t, getRes.Err)
	actual, err := getRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Post
	postOpts := &restclient.PostOpts{
		Params: params,
	}

	postRes := restclient.Post(c, url, postOpts)
	th.AssertNoErr(t, postRes.Err)
	actual, err = postRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Patch
	patchOpts := &restclient.PatchOpts{
		Params: params,
	}

	patchRes := restclient.Patch(c, url, patchOpts)
	th.AssertNoErr(t, patchRes.Err)
	actual, err = patchRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Put
	putOpts := &restclient.PutOpts{
		Params: params,
	}

	putRes := restclient.Put(c, url, putOpts)
	th.AssertNoErr(t, putRes.Err)
	actual, err = putRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Delete
	deleteRes := restclient.Delete(c, url, nil)
	th.AssertNoErr(t, deleteRes.Err)
	err = deleteRes.ExtractErr()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)
}

func TestNoContent(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	url := fmt.Sprintf("%s/route", th.Endpoint())

	c := new(gophercloud.ServiceClient)
	c.ProviderClient = new(gophercloud.ProviderClient)

	expected := map[string]interface{}(nil)

	// shared params
	params := map[string]interface{}{
		"bar": "baz",
	}

	// Get
	getRes := restclient.Get(c, url, nil)
	th.AssertNoErr(t, getRes.Err)
	actual, err := getRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Post
	postOpts := &restclient.PostOpts{
		Params: params,
	}

	postRes := restclient.Post(c, url, postOpts)
	th.AssertNoErr(t, postRes.Err)
	actual, err = postRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Patch
	patchOpts := &restclient.PatchOpts{
		Params: params,
	}

	patchRes := restclient.Patch(c, url, patchOpts)
	th.AssertNoErr(t, patchRes.Err)
	actual, err = patchRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Put
	putOpts := &restclient.PutOpts{
		Params: params,
	}

	putRes := restclient.Put(c, url, putOpts)
	th.AssertNoErr(t, putRes.Err)
	actual, err = putRes.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)

	// Delete
	deleteRes := restclient.Delete(c, url, nil)
	th.AssertNoErr(t, deleteRes.Err)
	err = deleteRes.ExtractErr()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)
}
