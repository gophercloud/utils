package helpers

import (
	"errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

// ProjectPurgeAll purges all the resources associated with a project.
// This includes: servers, snapshosts, volumes, floating IPs, routers, networks, sub-networks and security groups
func ProjectPurgeAll(computeClient *gophercloud.ServiceClient, storageClient *gophercloud.ServiceClient, networkClient *gophercloud.ServiceClient, projectID string) (err error) {
	err = ProjectPurgeCompute(computeClient, projectID)
	if err != nil {
		return err
	}
	err = ProjectPurgeStorage(storageClient, projectID)
	if err != nil {
		return err
	}
	err = ProjectPurgeNetwork(networkClient, projectID)
	if err != nil {
		return err
	}
	return nil
}

// ProjectPurgeCompute purges the Compute v2 resources associated with a project.
// This includes: servers
func ProjectPurgeCompute(computeClient *gophercloud.ServiceClient, projectID string) (err error) {
	// Delete servers
	listOpts := servers.ListOpts{
		AllTenants: true,
		TenantID: projectID,
	}

	allPages, err := servers.List(computeClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding servers for project: " + projectID)
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		return errors.New("Error extracting servers for project: " + projectID)
	}

	if len(allServers) > 0 {
		for _, server := range allServers {
			err = servers.Delete(computeClient, server.ID).ExtractErr()
			if err != nil {
				return errors.New("Error deleting server: " + server.Name + " from project: " + projectID)
			}
		}
	}

	return nil
}

// ProjectPurgeStorage purges the Blockstorage v3 resources associated with a project.
// This includes: snapshosts and volumes
func ProjectPurgeStorage(storageClient *gophercloud.ServiceClient, projectID string) (err error) {
	// Delete snapshots
	err = clearSnaphosts(projectID, storageClient)
	if err != nil {
		return err
	}
	// Delete volumes
	err = clearVolumes(projectID, storageClient)
	if err != nil {
		return err
	}

	return nil
}

// ProjectPurgeNetwork purges the Networking v2 resources associated with a project.
// This includes: floating IPs, routers, networks, sub-networks and security groups
func ProjectPurgeNetwork(networkClient *gophercloud.ServiceClient, projectID string) (err error) {
	// Delete floating IPs
	err = clearFloatings(projectID, networkClient)
	if err != nil {
		return err
	}
	// Delete ports
	err = clearPorts(projectID, networkClient)
	if err != nil {
		return err
	}
	// Delete routers
	err = clearRouters(projectID, networkClient)
	if err != nil {
		return err
	}
	// Delete networks
	err = clearNetworks(projectID, networkClient)
	if err != nil {
		return err
	}
	// Delete security groups
	err = clearSecGroups(projectID, networkClient)
	if err != nil {
		return err
	}

	return nil
}

func clearVolumes(projectID string, storageClient *gophercloud.ServiceClient) error {
	listOpts := volumes.ListOpts{
		AllTenants: true,
		TenantID:   projectID,
	}
	allPages, err := volumes.List(storageClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding volumes for project: " + projectID)
	}
	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return errors.New("Error extracting volumes for project: " + projectID)
	}
	if len(allVolumes) > 0 {
		deleteOpts := volumes.DeleteOpts{
			Cascade: true,
		}
		for _, volume := range allVolumes {
			err = volumes.Delete(storageClient, volume.ID, deleteOpts).ExtractErr()
			if err != nil {
				return errors.New("Error deleting volume: " + volume.Name + " from project: " + projectID)
			}
		}
	}

	return err
}

func clearSnaphosts(projectID string, storageClient *gophercloud.ServiceClient) error {
	listOpts := snapshots.ListOpts{
		AllTenants: true,
		TenantID:   projectID,
	}
	allPages, err := snapshots.List(storageClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding snapshots for project: " + projectID)
	}
	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		return errors.New("Error extracting snapshots for project: " + projectID)
	}
	if len(allSnapshots) > 0 {
		for _, snaphost := range allSnapshots {
			err = snapshots.Delete(storageClient, snaphost.ID).ExtractErr()
			if err != nil {
				return errors.New("Error deleting snaphost: " + snaphost.Name + " from project: " + projectID)
			}
		}
	}
	return nil
}

func clearFloatings(projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := floatingips.ListOpts{
		TenantID:   projectID,
	}
	allPages, err := floatingips.List(networkClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding floating IPs for project: " + projectID)
	}
	allFloatings, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return errors.New("Error extracting floating IPs for project: " + projectID)
	}
	if len(allFloatings) > 0 {
		for _, floating := range allFloatings {
			err = floatingips.Delete(networkClient, floating.ID).ExtractErr()
			if err != nil {
				return errors.New("Error deleting floating IP: " + floating.ID + " from project: " + projectID)
			}
		}
	}

	return nil
}

func clearPorts(projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := ports.ListOpts{
		TenantID:   projectID,
	}
	allPages, err := ports.List(networkClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding ports for project: " + projectID)
	}
	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return errors.New("Error extracting ports for project: " + projectID)
	}
	if len(allPorts) > 0 {
		for _, port := range allPorts {
			err = ports.Delete(networkClient, port.ID).ExtractErr()
			if err != nil {
				return errors.New("Error deleting port: " + port.ID + " from project: " + projectID)
			}
		}
	}

	return nil
}

func clearRouters(projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := routers.ListOpts{
		TenantID:   projectID,
	}
	allPages, err := routers.List(networkClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding routers for project: " + projectID)
	}
	allRouters, err := routers.ExtractRouters(allPages)
	if err != nil {
		return errors.New("Error extracting routers for project: " + projectID)
	}
	if len(allRouters) > 0 {
		for _, router := range allRouters {
			err = routers.Delete(networkClient, router.ID).ExtractErr()
			if err != nil {
				return errors.New("Error deleting router: " + router.Name + " from project: " + projectID)
			}
		}
	}

	return nil
}

func clearNetworks(projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := networks.ListOpts{
		TenantID:   projectID,
	}
	allPages, err := networks.List(networkClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding networks for project: " + projectID)
	}
	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return errors.New("Error extracting networks for project: " + projectID)
	}
	if len(allNetworks) > 0 {
		for _, network := range allNetworks {
			err = networks.Delete(networkClient, network.ID).ExtractErr()
			if err != nil {
				return errors.New("Error deleting network: " + network.Name + " from project: " +  projectID)
			}
		}
	}

	return nil
}

func clearSecGroups(projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := groups.ListOpts{
		TenantID:   projectID,
	}
	allPages, err := groups.List(networkClient, listOpts).AllPages()
	if err != nil {
		return errors.New("Error finding security groups for project: " + projectID)
	}
	allSecGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return errors.New("Error extracting security groups for project: " + projectID)
	}
	if len(allSecGroups) > 0 {
		for _, group := range allSecGroups {
			err = groups.Delete(networkClient, group.ID).ExtractErr()
			if err != nil {
				return errors.New("Error deleting security group: " + group.Name + " from project: " +  projectID)
			}
		}
	}

	return nil
}
