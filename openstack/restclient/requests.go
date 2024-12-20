package restclient

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
)

// GetOpts represents options used in the Get request.
type GetOpts struct {
	Headers map[string]string
	Query   map[string]interface{}
}

// Get performs a generic GET request to the specified URL.
func Get(client *gophercloud.ServiceClient, url string, opts *GetOpts) (r GetResult) {
	requestOpts := new(gophercloud.RequestOpts)

	if opts != nil {
		query, err := BuildQueryString(opts.Query)
		if err != nil {
			r.Err = err
			return
		}
		url += query.String()

		requestOpts.MoreHeaders = opts.Headers
	}

	// Allow a wide range of statuses
	requestOpts.OkCodes = []int{200, 201, 202, 204, 206}

	_, err := client.Get(url, &r.Body, requestOpts)
	if err != nil {
		if err.Error() != "EOF" {
			r.Err = err
			return
		}

		err = nil
		r.Body = nil
	}

	return
}

// PostOpts represents options used in a Post request.
type PostOpts struct {
	Headers map[string]string
	Params  map[string]interface{}
	Query   map[string]interface{}
}

// Post performs a generic POST request to the specified URL.
func Post(client *gophercloud.ServiceClient, url string, opts *PostOpts) (r PostResult) {
	var b map[string]interface{}
	requestOpts := new(gophercloud.RequestOpts)

	if opts != nil {
		query, err := BuildQueryString(opts.Query)
		if err != nil {
			r.Err = err
			return
		}
		url += query.String()

		b = opts.Params

		requestOpts.MoreHeaders = opts.Headers
	}

	// Allow a wide range of statuses
	requestOpts.OkCodes = []int{200, 201, 202, 204, 206}

	_, err := client.Post(url, &b, &r.Body, requestOpts)
	if err != nil {
		if err.Error() != "EOF" {
			r.Err = err
			return
		}

		err = nil
		r.Body = nil
	}

	return
}

// PutOpts represents options used in a Put request.
type PutOpts struct {
	Headers map[string]string
	Params  map[string]interface{}
	Query   map[string]interface{}
}

// Put performs a generic PUT request to the specified URL.
func Put(client *gophercloud.ServiceClient, url string, opts *PutOpts) (r PostResult) {
	var b map[string]interface{}
	requestOpts := new(gophercloud.RequestOpts)

	if opts != nil {
		query, err := BuildQueryString(opts.Query)
		if err != nil {
			r.Err = err
			return
		}
		url += query.String()

		b = opts.Params
	}

	// Allow a wide range of statuses
	requestOpts.OkCodes = []int{200, 201, 202, 204, 206}

	_, err := client.Put(url, &b, &r.Body, requestOpts)
	if err != nil {
		if err.Error() != "EOF" {
			r.Err = err
			return
		}

		err = nil
		r.Body = nil
	}

	return
}

// PatchOpts represents options used in a Patch request.
type PatchOpts struct {
	Headers map[string]string
	Params  map[string]interface{}
	Query   map[string]interface{}
}

// Patch performs a generic PATCH request to the specified URL.
func Patch(client *gophercloud.ServiceClient, url string, opts *PatchOpts) (r PatchResult) {
	var b map[string]interface{}
	requestOpts := new(gophercloud.RequestOpts)

	if opts != nil {
		query, err := BuildQueryString(opts.Query)
		if err != nil {
			r.Err = err
			return
		}
		url += query.String()

		b = opts.Params
	}

	// Allow a wide range of statuses
	requestOpts.OkCodes = []int{200, 201, 202, 204, 206}

	_, err := client.Patch(url, &b, &r.Body, requestOpts)
	if err != nil {
		if err.Error() != "EOF" {
			r.Err = err
			return
		}

		err = nil
		r.Body = nil
	}

	return
}

// DeleteOpts represents options used in a Delete request.
type DeleteOpts struct {
	Headers map[string]string
	Query   map[string]interface{}
}

// Delete performs a generic DELETE request to the specified URL.
func Delete(client *gophercloud.ServiceClient, url string, opts *DeleteOpts) (r DeleteResult) {
	requestOpts := new(gophercloud.RequestOpts)

	if opts != nil {
		query, err := BuildQueryString(opts.Query)
		if err != nil {
			r.Err = err
			return
		}
		url += query.String()

		requestOpts.MoreHeaders = opts.Headers
	}

	// Allow a wide range of statuses
	requestOpts.OkCodes = []int{200, 201, 202, 204, 206}

	_, err := client.Delete(url, requestOpts)
	if err != nil {
		if err.Error() != "EOF" {
			r.Err = err
			return
		}

		err = nil
		r.Body = nil
	}

	return
}

// BuildQueryString will take a map[string]interface and convert it
// to a URL encoded string. This is a watered-down version of Gophercloud's
// BuildQueryString.
func BuildQueryString(q map[string]interface{}) (*url.URL, error) {
	params := url.Values{}

	for key, value := range q {
		v := reflect.ValueOf(value)

		switch v.Kind() {
		case reflect.String:
			params.Add(key, v.String())
		case reflect.Int:
			params.Add(key, strconv.FormatInt(v.Int(), 10))
		case reflect.Bool:
			params.Add(key, strconv.FormatBool(v.Bool()))
		case reflect.Slice:
			switch v.Type().Elem() {
			case reflect.TypeOf(0):
				for i := 0; i < v.Len(); i++ {
					params.Add(key, strconv.FormatInt(v.Index(i).Int(), 10))
				}
			default:
				for i := 0; i < v.Len(); i++ {
					params.Add(key, v.Index(i).String())
				}
			}
		case reflect.Map:
			if v.Type().Key().Kind() == reflect.String && v.Type().Elem().Kind() == reflect.String {
				var s []string
				for _, k := range v.MapKeys() {
					value := v.MapIndex(k).String()
					s = append(s, fmt.Sprintf("'%s':'%s'", k.String(), value))
				}
				params.Add(key, fmt.Sprintf("{%s}", strings.Join(s, ", ")))
			}
		default:
			return nil, fmt.Errorf("Unsupported type: %s", v.Kind())
		}
	}

	return &url.URL{RawQuery: params.Encode()}, nil
}
