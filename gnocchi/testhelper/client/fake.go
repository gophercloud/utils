package client

import (
	"github.com/gophercloud/gophercloud/v2"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// TokenID is a fake Identity service token.
const TokenID = client.TokenID

// ServiceClient returns a generic service client for use in tests.
func ServiceClient(fakeServer th.FakeServer) *gophercloud.ServiceClient {
	sc := client.ServiceClient(fakeServer)
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc
}
