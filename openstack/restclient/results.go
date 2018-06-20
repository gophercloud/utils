package restclient

import (
	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (map[string]interface{}, error) {
	var s map[string]interface{}
	err := r.ExtractInto(&s)
	return s, err
}

type GetResult struct {
	commonResult
}

type PostResult struct {
	commonResult
}

type PatchResult struct {
	commonResult
}

type PutResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
