package gnocchi

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

func initClientOpts(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, clientType string) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := client.EndpointLocator(ctx, eo)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	sc.Endpoint = url
	sc.Type = clientType
	return sc, nil
}

// NewGnocchiV1 creates a ServiceClient that may be used with the v1 Gnocchi package.
func NewGnocchiV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "metric")
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc, err
}
