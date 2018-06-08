// +build acceptance metric archivepolicies

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/utils/acceptance/clients"
	"github.com/gophercloud/utils/gnocchi/metric/v1/archivepolicies"
)

func TestArchivePoliciesCRUD(t *testing.T) {
	client, err := clients.NewGnocchiV1Client()
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi client: %v", err)
	}

	archivePolicy, err := CreateArchivePolicy(t, client)
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi archive policy: %v", err)
	}
	// defer DeleteArchivePolicy(t, client, archivePolicy.ID)

	tools.PrintResource(t, archivePolicy)
}

func TestArchivePoliciesList(t *testing.T) {
	client, err := clients.NewGnocchiV1Client()
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi client: %v", err)
	}

	allPages, err := archivepolicies.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list archive policies: %v", err)
	}

	allArchivePolicies, err := archivepolicies.ExtractArchivePolicies(allPages)
	if err != nil {
		t.Fatalf("Unable to extract archive policies: %v", err)
	}

	for _, archivePolicy := range allArchivePolicies {
		tools.PrintResource(t, archivePolicy)
	}
}
