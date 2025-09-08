package helpers

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/portforwarding"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
)

type ProjectPurgeOpts struct {
	ComputePurgeOpts *ComputePurgeOpts
	StoragePurgeOpts *StoragePurgeOpts
	NetworkPurgeOpts *NetworkPurgeOpts
}

type ComputePurgeOpts struct {
	// Client is a reference to a specific compute service client.
	Client *gophercloud.ServiceClient
}

type StoragePurgeOpts struct {
	// Client is a reference to a specific storage service client.
	Client *gophercloud.ServiceClient
}

type NetworkPurgeOpts struct {
	// Client is a reference to a specific networking service client.
	Client *gophercloud.ServiceClient
}

// ProjectPurgeAll purges all the resources associated with a project.
// This includes: servers, snapshosts, volumes, floating IPs, routers, networks, sub-networks and security groups
func ProjectPurgeAll(ctx context.Context, projectID string, purgeOpts ProjectPurgeOpts) (err error) {
	if purgeOpts.ComputePurgeOpts != nil {
		err = ProjectPurgeCompute(ctx, projectID, *purgeOpts.ComputePurgeOpts)
		if err != nil {
			return err
		}
	}
	if purgeOpts.StoragePurgeOpts != nil {
		err = ProjectPurgeStorage(ctx, projectID, *purgeOpts.StoragePurgeOpts)
		if err != nil {
			return err
		}
	}
	if purgeOpts.NetworkPurgeOpts != nil {
		err = ProjectPurgeNetwork(ctx, projectID, *purgeOpts.NetworkPurgeOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProjectPurgeCompute purges the Compute v2 resources associated with a project.
// This includes: servers
func ProjectPurgeCompute(ctx context.Context, projectID string, purgeOpts ComputePurgeOpts) (err error) {
	// Delete servers
	listOpts := servers.ListOpts{
		AllTenants: true,
		TenantID:   projectID,
	}

	allPages, err := servers.List(purgeOpts.Client, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding servers for project: %s", projectID)
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		return fmt.Errorf("error extracting servers for project: %s", projectID)
	}

	if len(allServers) > 0 {
		for _, server := range allServers {
			err = servers.Delete(ctx, purgeOpts.Client, server.ID).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting server: %s from project: %s", server.Name, projectID)
			}
		}
	}

	return nil
}

// ProjectPurgeStorage purges the Blockstorage v3 resources associated with a project.
// This includes: snapshosts and volumes
func ProjectPurgeStorage(ctx context.Context, projectID string, purgeOpts StoragePurgeOpts) (err error) {
	// Delete snapshots
	err = clearBlockStorageSnaphosts(ctx, projectID, purgeOpts.Client)
	if err != nil {
		return err
	}
	// Delete volumes
	err = clearBlockStorageVolumes(ctx, projectID, purgeOpts.Client)
	if err != nil {
		return err
	}

	return nil
}

// ProjectPurgeNetwork purges the Networking v2 resources associated with a project.
// This includes: floating IPs, routers, networks, sub-networks and security groups
func ProjectPurgeNetwork(ctx context.Context, projectID string, purgeOpts NetworkPurgeOpts) (err error) {
	// Delete floating IPs
	err = clearNetworkingFloatingIPs(ctx, projectID, purgeOpts.Client)
	if err != nil {
		return err
	}
	// Delete ports
	err = clearNetworkingPorts(ctx, projectID, purgeOpts.Client)
	if err != nil {
		return err
	}
	// Delete routers
	err = clearNetworkingRouters(ctx, projectID, purgeOpts.Client)
	if err != nil {
		return err
	}
	// Delete networks
	err = clearNetworkingNetworks(ctx, projectID, purgeOpts.Client)
	if err != nil {
		return err
	}
	// Delete security groups
	err = clearNetworkingSecurityGroups(ctx, projectID, purgeOpts.Client)
	if err != nil {
		return err
	}

	return nil
}

func clearBlockStorageVolumes(ctx context.Context, projectID string, storageClient *gophercloud.ServiceClient) error {
	listOpts := volumes.ListOpts{
		AllTenants: true,
		TenantID:   projectID,
	}
	allPages, err := volumes.List(storageClient, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding volumes for project: %s", projectID)
	}
	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return fmt.Errorf("error extracting volumes for project: %s", projectID)
	}
	if len(allVolumes) > 0 {
		deleteOpts := volumes.DeleteOpts{
			Cascade: true,
		}
		for _, volume := range allVolumes {
			err = volumes.Delete(ctx, storageClient, volume.ID, deleteOpts).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting volume: %s from project: %s", volume.Name, projectID)
			}
		}
	}

	return err
}

func clearBlockStorageSnaphosts(ctx context.Context, projectID string, storageClient *gophercloud.ServiceClient) error {
	listOpts := snapshots.ListOpts{
		AllTenants: true,
		TenantID:   projectID,
	}
	allPages, err := snapshots.List(storageClient, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding snapshots for project: %s", projectID)
	}
	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		return fmt.Errorf("error extracting snapshots for project: %s", projectID)
	}
	if len(allSnapshots) > 0 {
		for _, snaphost := range allSnapshots {
			err = snapshots.Delete(ctx, storageClient, snaphost.ID).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting snaphost: %s from project: %s", snaphost.Name, projectID)
			}
		}
	}
	return nil
}

func clearPortforwarding(ctx context.Context, networkClient *gophercloud.ServiceClient, fipID string, projectID string) error {
	allPages, err := portforwarding.List(networkClient, portforwarding.ListOpts{}, fipID).AllPages(ctx)
	if err != nil {
		return err
	}

	allPFs, err := portforwarding.ExtractPortForwardings(allPages)
	if err != nil {
		return err
	}

	for _, pf := range allPFs {
		err := portforwarding.Delete(ctx, networkClient, fipID, pf.ID).ExtractErr()
		if err != nil {
			return fmt.Errorf("error deleting floating IP port forwarding: %s from project %s", pf.ID, projectID)
		}
	}

	return nil
}

func clearNetworkingFloatingIPs(ctx context.Context, projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := floatingips.ListOpts{
		TenantID: projectID,
	}
	allPages, err := floatingips.List(networkClient, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding floating IPs for project: %s", projectID)
	}
	allFloatings, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return fmt.Errorf("error extracting floating IPs for project: %s", projectID)
	}
	if len(allFloatings) > 0 {
		for _, floating := range allFloatings {
			// Clear all portforwarding settings otherwise the floating IP can't be deleted
			err = clearPortforwarding(ctx, networkClient, floating.ID, projectID)
			if err != nil {
				return err
			}

			err = floatingips.Delete(ctx, networkClient, floating.ID).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting floating IP: %s from project: %s", floating.ID, projectID)
			}
		}
	}

	return nil
}

func clearNetworkingPorts(ctx context.Context, projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := ports.ListOpts{
		TenantID: projectID,
	}

	allPages, err := ports.List(networkClient, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding ports for project: %s", projectID)
	}
	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return fmt.Errorf("error extracting ports for project: %s", projectID)
	}
	if len(allPorts) > 0 {
		for _, port := range allPorts {
			if port.DeviceOwner == "network:ha_router_replicated_interface" {
				continue
			}

			err = ports.Delete(ctx, networkClient, port.ID).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting port: %s from project: %s", port.ID, projectID)
			}
		}
	}

	return nil
}

// We need all subnets to disassociate the router from the subnet
func getAllSubnets(ctx context.Context, projectID string, networkClient *gophercloud.ServiceClient) ([]string, error) {
	subnets := make([]string, 0)
	listOpts := networks.ListOpts{
		TenantID: projectID,
	}

	allPages, err := networks.List(networkClient, listOpts).AllPages(ctx)
	if err != nil {
		return subnets, fmt.Errorf("error finding networks for project: %s", projectID)
	}
	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return subnets, fmt.Errorf("error extracting networks for project: %s", projectID)
	}
	if len(allNetworks) > 0 {
		for _, network := range allNetworks {
			subnets = append(subnets, network.Subnets...)
		}
	}

	return subnets, nil
}

func clearAllRouterInterfaces(ctx context.Context, routerID string, subnets []string, networkClient *gophercloud.ServiceClient) error {
	for _, subnet := range subnets {
		intOpts := routers.RemoveInterfaceOpts{
			SubnetID: subnet,
		}

		_, err := routers.RemoveInterface(ctx, networkClient, routerID, intOpts).Extract()
		if err != nil {
			return err
		}
	}

	return nil
}

func clearNetworkingRouters(ctx context.Context, projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := routers.ListOpts{
		TenantID: projectID,
	}
	allPages, err := routers.List(networkClient, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding routers for project: %s", projectID)
	}
	allRouters, err := routers.ExtractRouters(allPages)
	if err != nil {
		return fmt.Errorf("error extracting routers for project: %s", projectID)
	}

	subnets, err := getAllSubnets(ctx, projectID, networkClient)
	if err != nil {
		return fmt.Errorf("error fetching subnets project: %s", projectID)
	}

	if len(allRouters) > 0 {
		for _, router := range allRouters {
			err = clearAllRouterInterfaces(ctx, router.ID, subnets, networkClient)
			if err != nil {
				return err
			}

			routes := []routers.Route{}
			// Clear all routes
			updateOpts := routers.UpdateOpts{
				Routes: &routes,
			}

			_, err := routers.Update(ctx, networkClient, router.ID, updateOpts).Extract()
			if err != nil {
				return fmt.Errorf("error deleting router: %s from project: %s", router.Name, projectID)
			}

			err = routers.Delete(ctx, networkClient, router.ID).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting router: %s from project: %s", router.Name, projectID)
			}
		}
	}

	return nil
}

func clearNetworkingNetworks(ctx context.Context, projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := networks.ListOpts{
		TenantID: projectID,
	}

	allPages, err := networks.List(networkClient, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding networks for project: %s", projectID)
	}
	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return fmt.Errorf("error extracting networks for project: %s", projectID)
	}
	if len(allNetworks) > 0 {
		for _, network := range allNetworks {
			err = networks.Delete(ctx, networkClient, network.ID).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting network: %s from project: %s", network.Name, projectID)
			}
		}
	}

	return nil
}

func clearNetworkingSecurityGroups(ctx context.Context, projectID string, networkClient *gophercloud.ServiceClient) error {
	listOpts := groups.ListOpts{
		TenantID: projectID,
	}
	allPages, err := groups.List(networkClient, listOpts).AllPages(ctx)
	if err != nil {
		return fmt.Errorf("error finding security groups for project: %s", projectID)
	}
	allSecGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return fmt.Errorf("error extracting security groups for project: %s", projectID)
	}
	if len(allSecGroups) > 0 {
		for _, group := range allSecGroups {
			err = groups.Delete(ctx, networkClient, group.ID).ExtractErr()
			if err != nil {
				return fmt.Errorf("error deleting security group: %s from project: %s", group.Name, projectID)
			}
		}
	}

	return nil
}
