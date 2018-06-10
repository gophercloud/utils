package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/utils/gnocchi/metric/v1/archivepolicies"
)

// CreateArchivePolicy will create a Gnocchi archive policy. An error will be returned if the
// archive policy could not be created.
func CreateArchivePolicy(t *testing.T, client *gophercloud.ServiceClient) (*archivepolicies.ArchivePolicy, error) {
	policyName := tools.RandomString("TESTACCT-", 8)
	createOpts := archivepolicies.CreateOpts{
		Name: policyName,
		AggregationMethods: []string{
			"mean",
			"sum",
		},
		Definition: []archivepolicies.ArchivePolicyDefinitionOpts{
			{
				Granularity: "1:00:00",
				TimeSpan:    "30 days, 0:00:00",
			},
			{
				Granularity: "24:00:00",
				TimeSpan:    "90 days, 0:00:00",
			},
		},
	}

	t.Logf("Attempting to create a Gnocchi archive policy")
	archivePolicy, err := archivepolicies.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the Gnocchi archive policy.")
	return archivePolicy, nil
}
