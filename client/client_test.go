package client

import (
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestFormatHeaders(t *testing.T) {
	headers := http.Header{
		"X-Auth-Token": []string{"token"},
		"User-Agent":   []string{"Terraform/x.x.x", "Gophercloud/y.y.y"},
	}

	expected := "User-Agent: Terraform/x.x.x Gophercloud/y.y.y\nX-Auth-Token: ***"
	rt := RoundTripper{}
	actual := rt.formatHeaders(headers, "\n")

	th.AssertEquals(t, expected, actual)
}
