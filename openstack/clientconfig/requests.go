package clientconfig

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	"gopkg.in/yaml.v2"
)

// ClientOpts represents options to customize the way a client is
// configured.
type ClientOpts struct {
	// Cloud is the cloud entry in clouds.yaml to use.
	Cloud string

	// EnvPrefix allows a custom environment variable prefix to be used.
	EnvPrefix string

	// AuthType specifies the type of authentication to use.
	// By default, this is "password".
	AuthType string

	// Auth defines the authentication information needed to
	// authenticate to a cloud when clouds.yaml isn't used.
	Auth *Auth
}

// LoadYAML will load a clouds.yaml file and return the full config.
func LoadYAML() (map[string]Cloud, error) {
	content, err := findAndReadYAML()
	if err != nil {
		return nil, err
	}

	var clouds Clouds
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return clouds.Clouds, nil
}

// GetCloudFromYAML will return a cloud entry from a clouds.yaml file.
func GetCloudFromYAML(opts *ClientOpts) (*Cloud, error) {
	clouds, err := LoadYAML()
	if err != nil {
		return nil, fmt.Errorf("unable to load clouds.yaml: %s", err)
	}

	// Determine which cloud to use.
	// First see if a cloud name was explicitly set in opts.
	var cloudName string
	if opts != nil && opts.Cloud != "" {
		cloudName = opts.Cloud
	}

	// Next see if a cloud name was specified as an environment variable.
	// This is supposed to override an explicit opts setting.
	envPrefix := "OS_"
	if opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	if v := os.Getenv(envPrefix + "CLOUD"); v != "" {
		cloudName = v
	}

	var cloud *Cloud
	if cloudName != "" {
		v, ok := clouds[cloudName]
		if !ok {
			return nil, fmt.Errorf("cloud %s does not exist in clouds.yaml", cloudName)
		}
		cloud = &v
	}

	// If a cloud was not specified, and clouds only contains
	// a single entry, use that entry.
	if cloudName == "" && len(clouds) == 1 {
		for _, v := range clouds {
			cloud = &v
		}
	}

	if cloud == nil {
		return nil, fmt.Errorf("Unable to determine a valid entry in clouds.yaml")
	}

	return cloud, nil
}

// AuthOptions creates a gophercloud.AuthOptions structure with the
// settings found in a specific cloud entry of a clouds.yaml file or
// based on authentication settings given in ClientOpts.
//
// This attempts to be a single point of entry for all OpenStack authentication.
//
// See http://docs.openstack.org/developer/os-client-config and
// https://github.com/openstack/os-client-config/blob/master/os_client_config/config.py.
func AuthOptions(opts *ClientOpts) (*gophercloud.AuthOptions, error) {
	cloud := new(Cloud)

	// If no opts were passed in, create an empty ClientOpts.
	if opts == nil {
		opts = new(ClientOpts)
	}

	// Determine if a clouds.yaml entry should be retrieved.
	// Start by figuring out the cloud name.
	// First check if one was explicitly specified in opts.
	var cloudName string
	if opts.Cloud != "" {
		cloudName = opts.Cloud
	}

	// Next see if a cloud name was specified as an environment variable.
	envPrefix := "OS_"
	if opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	if v := os.Getenv(envPrefix + "CLOUD"); v != "" {
		cloudName = v
	}

	// If a cloud name was determined, try to look it up in clouds.yaml.
	if cloudName != "" {
		// Get the requested cloud.
		var err error
		cloud, err = GetCloudFromYAML(opts)
		if err != nil {
			return nil, err
		}
	}

	// If cloud.Auth is nil, then no cloud was specified.
	if cloud.Auth == nil {
		// If opts.Auth is not nil, then try using the auth settings from it.
		if opts.Auth != nil {
			cloud.Auth = opts.Auth
		}

		// If cloud.Auth is still nil, then set it to an empty Auth struct
		// and rely on environment variables to do the authentication.
		if cloud.Auth == nil {
			cloud.Auth = new(Auth)
		}
	}

	identityAPI := determineIdentityAPI(cloud, opts)
	switch identityAPI {
	case "2.0":
		return v2auth(cloud, opts)
	case "3":
		return v3auth(cloud, opts)
	}

	return nil, fmt.Errorf("Unable to build AuthOptions")
}

func determineIdentityAPI(cloud *Cloud, opts *ClientOpts) string {
	var identityAPI string
	if cloud.IdentityAPIVersion != "" {
		identityAPI = cloud.IdentityAPIVersion
	}

	envPrefix := "OS_"
	if opts != nil && opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	if v := os.Getenv(envPrefix + "IDENTITY_API_VERSION"); v != "" {
		identityAPI = v
	}

	if identityAPI == "" {
		if cloud.Auth != nil {
			if strings.Contains(cloud.Auth.AuthURL, "v2.0") {
				identityAPI = "2.0"
			}

			if strings.Contains(cloud.Auth.AuthURL, "v3") {
				identityAPI = "3"
			}
		}
	}

	if identityAPI == "" {
		switch cloud.AuthType {
		case "v2password":
			identityAPI = "2.0"
		case "v2token":
			identityAPI = "2.0"
		case "v3password":
			identityAPI = "3"
		case "v3token":
			identityAPI = "3"
		}
	}

	// If an Identity API version could not be determined,
	// default to v3.
	if identityAPI == "" {
		identityAPI = "3"
	}

	return identityAPI
}

// v2auth creates a v2-compatible gophercloud.AuthOptions struct.
func v2auth(cloud *Cloud, opts *ClientOpts) (*gophercloud.AuthOptions, error) {
	// Environment variable overrides.
	envPrefix := "OS_"
	if opts != nil && opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	if v := os.Getenv(envPrefix + "AUTH_URL"); v != "" {
		cloud.Auth.AuthURL = v
	}

	if v := os.Getenv(envPrefix + "TOKEN"); v != "" {
		cloud.Auth.Token = v
	}

	if v := os.Getenv(envPrefix + "AUTH_TOKEN"); v != "" {
		cloud.Auth.Token = v
	}

	if v := os.Getenv(envPrefix + "USERNAME"); v != "" {
		cloud.Auth.Username = v
	}

	if v := os.Getenv(envPrefix + "PASSWORD"); v != "" {
		cloud.Auth.Password = v
	}

	if v := os.Getenv(envPrefix + "TENANT_ID"); v != "" {
		cloud.Auth.ProjectID = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_ID"); v != "" {
		cloud.Auth.ProjectID = v
	}

	if v := os.Getenv(envPrefix + "TENANT_NAME"); v != "" {
		cloud.Auth.ProjectName = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_NAME"); v != "" {
		cloud.Auth.ProjectName = v
	}

	ao := &gophercloud.AuthOptions{
		IdentityEndpoint: cloud.Auth.AuthURL,
		TokenID:          cloud.Auth.Token,
		Username:         cloud.Auth.Username,
		Password:         cloud.Auth.Password,
		TenantID:         cloud.Auth.ProjectID,
		TenantName:       cloud.Auth.ProjectName,
	}

	return ao, nil
}

// v3auth creates a v3-compatible gophercloud.AuthOptions struct.
func v3auth(cloud *Cloud, opts *ClientOpts) (*gophercloud.AuthOptions, error) {
	// Environment variable overrides.
	envPrefix := "OS_"
	if opts != nil && opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	if v := os.Getenv(envPrefix + "AUTH_URL"); v != "" {
		cloud.Auth.AuthURL = v
	}

	if v := os.Getenv(envPrefix + "TOKEN"); v != "" {
		cloud.Auth.Token = v
	}

	if v := os.Getenv(envPrefix + "AUTH_TOKEN"); v != "" {
		cloud.Auth.Token = v
	}

	if v := os.Getenv(envPrefix + "USERNAME"); v != "" {
		cloud.Auth.Username = v
	}

	if v := os.Getenv(envPrefix + "USER_ID"); v != "" {
		cloud.Auth.UserID = v
	}

	if v := os.Getenv(envPrefix + "PASSWORD"); v != "" {
		cloud.Auth.Password = v
	}

	if v := os.Getenv(envPrefix + "TENANT_ID"); v != "" {
		cloud.Auth.ProjectID = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_ID"); v != "" {
		cloud.Auth.ProjectID = v
	}

	if v := os.Getenv(envPrefix + "TENANT_NAME"); v != "" {
		cloud.Auth.ProjectName = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_NAME"); v != "" {
		cloud.Auth.ProjectName = v
	}

	if v := os.Getenv(envPrefix + "DOMAIN_ID"); v != "" {
		cloud.Auth.DomainID = v
	}

	if v := os.Getenv(envPrefix + "DOMAIN_NAME"); v != "" {
		cloud.Auth.DomainName = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_DOMAIN_ID"); v != "" {
		cloud.Auth.ProjectDomainID = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_DOMAIN_NAME"); v != "" {
		cloud.Auth.ProjectDomainName = v
	}

	if v := os.Getenv(envPrefix + "USER_DOMAIN_ID"); v != "" {
		cloud.Auth.UserDomainID = v
	}

	if v := os.Getenv(envPrefix + "USER_DOMAIN_NAME"); v != "" {
		cloud.Auth.UserDomainName = v
	}

	// Build a scope and try to do it correctly.
	// https://github.com/openstack/os-client-config/blob/master/os_client_config/config.py#L595
	scope := new(gophercloud.AuthScope)

	if !isProjectScoped(cloud.Auth) {
		if cloud.Auth.DomainID != "" {
			scope.DomainID = cloud.Auth.DomainID
		} else if cloud.Auth.DomainName != "" {
			scope.DomainName = cloud.Auth.DomainName
		}
	} else {
		// If Domain* is set, but UserDomain* or ProjectDomain* aren't,
		// then use Domain* as the default setting.
		cloud = setDomainIfNeeded(cloud)

		if cloud.Auth.ProjectID != "" {
			scope.ProjectID = cloud.Auth.ProjectID
		} else {
			scope.ProjectName = cloud.Auth.ProjectName
			scope.DomainID = cloud.Auth.ProjectDomainID
			scope.DomainName = cloud.Auth.ProjectDomainName
		}
	}

	ao := &gophercloud.AuthOptions{
		Scope:            scope,
		IdentityEndpoint: cloud.Auth.AuthURL,
		TokenID:          cloud.Auth.Token,
		Username:         cloud.Auth.Username,
		UserID:           cloud.Auth.UserID,
		Password:         cloud.Auth.Password,
		TenantID:         cloud.Auth.ProjectID,
		TenantName:       cloud.Auth.ProjectName,
		DomainID:         cloud.Auth.UserDomainID,
		DomainName:       cloud.Auth.UserDomainName,
	}

	// If an auth_type of "token" was specified, then make sure
	// Gophercloud properly authenticates with a token. This involves
	// unsetting a few other auth options. The reason this is done
	// here is to wait until all auth settings (both in clouds.yaml
	// and via environment variables) are set and then unset them.
	if strings.Contains(cloud.AuthType, "token") || ao.TokenID != "" {
		ao.Username = ""
		ao.Password = ""
		ao.UserID = ""
		ao.DomainID = ""
		ao.DomainName = ""
	}

	// Check for absolute minimum requirements.
	if ao.IdentityEndpoint == "" {
		err := gophercloud.ErrMissingInput{Argument: "authURL"}
		return nil, err
	}

	return ao, nil
}

// AuthenticatedClient is a convenience function to get a new provider client
// based on a clouds.yaml entry.
func AuthenticatedClient(opts *ClientOpts) (*gophercloud.ProviderClient, error) {
	ao, err := AuthOptions(opts)
	if err != nil {
		return nil, err
	}

	return openstack.AuthenticatedClient(*ao)
}

// NewServiceClient is a convenience function to get a new service client.
func NewServiceClient(service string, opts *ClientOpts) (*gophercloud.ServiceClient, error) {
	var cloud *Cloud
	if opts.Cloud != "" {
		// Get the requested cloud.
		var err error
		cloud, err = GetCloudFromYAML(opts)
		if err != nil {
			return nil, err
		}
	}

	if cloud == nil {
		cloud.Auth = opts.Auth
	}

	// Environment variable overrides.
	envPrefix := "OS_"
	if opts != nil && opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	// Get a Provider Client
	pClient, err := AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	// Determine the region to use.
	var region string
	if v := cloud.RegionName; v != "" {
		region = cloud.RegionName
	}

	if v := os.Getenv(envPrefix + "REGION_NAME"); v != "" {
		region = v
	}

	eo := gophercloud.EndpointOpts{
		Region: region,
	}

	switch service {
	case "compute":
		return openstack.NewComputeV2(pClient, eo)
	case "database":
		return openstack.NewDBV1(pClient, eo)
	case "dns":
		return openstack.NewDNSV2(pClient, eo)
	case "identity":
		identityVersion := "3"
		if v := cloud.IdentityAPIVersion; v != "" {
			identityVersion = v
		}

		switch identityVersion {
		case "v2", "2", "2.0":
			return openstack.NewIdentityV2(pClient, eo)
		case "v3", "3":
			return openstack.NewIdentityV3(pClient, eo)
		default:
			return nil, fmt.Errorf("invalid identity API version")
		}
	case "image":
		return openstack.NewImageServiceV2(pClient, eo)
	case "network":
		return openstack.NewNetworkV2(pClient, eo)
	case "object-store":
		return openstack.NewObjectStorageV1(pClient, eo)
	case "orchestration":
		return openstack.NewOrchestrationV1(pClient, eo)
	case "sharev2":
		return openstack.NewSharedFileSystemV2(pClient, eo)
	case "volume":
		volumeVersion := "2"
		if v := cloud.VolumeAPIVersion; v != "" {
			volumeVersion = v
		}

		switch volumeVersion {
		case "v1", "1":
			return openstack.NewBlockStorageV1(pClient, eo)
		case "v2", "2":
			return openstack.NewBlockStorageV2(pClient, eo)
		default:
			return nil, fmt.Errorf("invalid volume API version")
		}
	}

	return nil, fmt.Errorf("unable to create a service client for %s", service)
}
