package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/swauth"
	osClient "github.com/gophercloud/utils/v2/client"
	"github.com/gophercloud/utils/v2/internal"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/gophercloud/utils/v2/terraform/mutexkv"
)

type Config struct {
	CACertFile                  string
	ClientCertFile              string
	ClientKeyFile               string
	Cloud                       string
	DefaultDomain               string
	DomainID                    string
	DomainName                  string
	EndpointOverrides           map[string]interface{}
	EndpointType                string
	IdentityEndpoint            string
	Insecure                    *bool
	Password                    string
	ProjectDomainName           string
	ProjectDomainID             string
	Region                      string
	Swauth                      bool
	TenantID                    string
	TenantName                  string
	Token                       string
	UserDomainName              string
	UserDomainID                string
	Username                    string
	UserID                      string
	ApplicationCredentialID     string
	ApplicationCredentialName   string
	ApplicationCredentialSecret string
	UseOctavia                  bool
	MaxRetries                  int
	DisableNoCacheHeader        bool
	Context                     context.Context

	DelayedAuth   bool
	AllowReauth   bool
	OsClient      *gophercloud.ProviderClient
	AuthOpts      *gophercloud.AuthOptions
	authenticated bool
	authFailed    error
	swClient      *gophercloud.ServiceClient
	swAuthFailed  error

	TerraformVersion string
	SDKVersion       string
	EnableLogger     bool

	*mutexkv.MutexKV
}

// LoadAndValidate performs the authentication and initial configuration
// of an OpenStack Provider Client. This sets up the HTTP client and
// authenticates to an OpenStack cloud.
//
// Individual Service Clients are created later in this file.
func (c *Config) LoadAndValidate(ctx context.Context) error {
	// Make sure at least one of auth_url or cloud was specified.
	if c.IdentityEndpoint == "" && c.Cloud == "" {
		return fmt.Errorf("One of 'auth_url' or 'cloud' must be specified")
	}

	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	clientOpts := new(clientconfig.ClientOpts)

	// If a cloud entry was given, base AuthOptions on a clouds.yaml file.
	if c.Cloud != "" {
		clientOpts.Cloud = c.Cloud

		// Passing region allows GetCloudFromYAML to apply per-region overrides
		clientOpts.RegionName = c.Region

		cloud, err := clientconfig.GetCloudFromYAML(clientOpts)
		if err != nil {
			return err
		}

		if c.Region == "" && cloud.RegionName != "" {
			c.Region = cloud.RegionName
		}

		if c.CACertFile == "" && cloud.CACertFile != "" {
			c.CACertFile = cloud.CACertFile
		}

		if c.ClientCertFile == "" && cloud.ClientCertFile != "" {
			c.ClientCertFile = cloud.ClientCertFile
		}

		if c.ClientKeyFile == "" && cloud.ClientKeyFile != "" {
			c.ClientKeyFile = cloud.ClientKeyFile
		}

		if c.Insecure == nil && cloud.Verify != nil {
			v := (!*cloud.Verify)
			c.Insecure = &v
		}
	} else {
		authInfo := &clientconfig.AuthInfo{
			AuthURL:                     c.IdentityEndpoint,
			DefaultDomain:               c.DefaultDomain,
			DomainID:                    c.DomainID,
			DomainName:                  c.DomainName,
			Password:                    c.Password,
			ProjectDomainID:             c.ProjectDomainID,
			ProjectDomainName:           c.ProjectDomainName,
			ProjectID:                   c.TenantID,
			ProjectName:                 c.TenantName,
			Token:                       c.Token,
			UserDomainID:                c.UserDomainID,
			UserDomainName:              c.UserDomainName,
			Username:                    c.Username,
			UserID:                      c.UserID,
			ApplicationCredentialID:     c.ApplicationCredentialID,
			ApplicationCredentialName:   c.ApplicationCredentialName,
			ApplicationCredentialSecret: c.ApplicationCredentialSecret,
		}

		// Define System Scope if enabled
		if c.AuthOpts.Scope.System {
			authInfo.SystemScope = "true"
		}

		clientOpts.AuthInfo = authInfo
	}

	ao, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] OpenStack allowReauth: %t", c.AllowReauth)
	ao.AllowReauth = c.AllowReauth

	client, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	// Set UserAgent
	client.UserAgent.Prepend(terraformUserAgent(c.TerraformVersion, c.SDKVersion))

	config, err := internal.PrepareTLSConfig(c.CACertFile, c.ClientCertFile, c.ClientKeyFile, c.Insecure)
	if err != nil {
		return err
	}

	c.EnableLogger = enableLogging(c.EnableLogger)
	var logger osClient.Logger
	if c.EnableLogger {
		logger = &osClient.DefaultLogger{}
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &osClient.RoundTripper{
			Rt:         transport,
			MaxRetries: c.MaxRetries,
			Logger:     logger,
		},
	}

	if !c.DisableNoCacheHeader {
		extraHeaders := map[string][]string{
			"Cache-Control": {"no-cache"},
		}
		client.HTTPClient.Transport.(*osClient.RoundTripper).SetHeaders(extraHeaders)
	}

	if c.MaxRetries > 0 {
		client.MaxBackoffRetries = uint(c.MaxRetries)
		client.RetryBackoffFunc = osClient.RetryBackoffFunc(logger)
	}

	if !c.DelayedAuth && !c.Swauth {
		err = openstack.Authenticate(ctx, client, *ao)
		if err != nil {
			return err
		}
	}

	c.AuthOpts = ao
	c.OsClient = client

	return nil
}

func (c *Config) Authenticate(ctx context.Context) error {
	if !c.DelayedAuth {
		return nil
	}

	c.MutexKV.Lock("auth")
	defer c.MutexKV.Unlock("auth")

	if c.authFailed != nil {
		return c.authFailed
	}

	if !c.authenticated {
		if err := openstack.Authenticate(ctx, c.OsClient, *c.AuthOpts); err != nil {
			c.authFailed = err
			return err
		}
		c.authenticated = true
	}

	return nil
}

// DetermineEndpoint is a helper method to determine if the user wants to
// override an endpoint returned from the catalog.
func (c *Config) DetermineEndpoint(client *gophercloud.ServiceClient, eo gophercloud.EndpointOpts, service string) (*gophercloud.ServiceClient, error) {
	v, ok := c.EndpointOverrides[service]
	if !ok {
		return client, nil
	}
	val, ok := v.(string)
	if !ok || val == "" {
		return client, nil
	}

	// overriden endpoint is a URL
	if u, err := url.Parse(val); err == nil && u.Scheme != "" && u.Host != "" {
		eo.ApplyDefaults(service)
		client.ProviderClient = c.OsClient
		client.Endpoint = val
		client.ResourceBase = ""
		client.Type = service
		log.Printf("[DEBUG] OpenStack Endpoint for %s: %s", service, val)
		return client, nil
	}

	// overriden endpoint is a new service type
	eo.ApplyDefaults(val)
	url, err := c.OsClient.EndpointLocator(eo)
	if err != nil {
		log.Printf("[DEBUG] Cannot set a new OpenStack Endpoint %s alias: %v", val, err)
		return client, err
	}
	client.ProviderClient = c.OsClient
	client.Endpoint = url
	client.Type = val

	log.Printf("[DEBUG] OpenStack Endpoint for %s alias: %s", val, url)
	return client, nil
}

// DetermineRegion is a helper method to determine the region based on
// the user's settings.
func (c *Config) DetermineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] OpenStack Region is: %s", region)
	return region
}

// The following methods assist with the creation of individual Service Clients
// which interact with the various OpenStack services.

type commonCommonServiceClientInitFunc func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)

func (c *Config) CommonServiceClientInit(ctx context.Context, newClient commonCommonServiceClientInitFunc, region, service string) (*gophercloud.ServiceClient, error) {
	if err := c.Authenticate(ctx); err != nil {
		return nil, err
	}

	eo := gophercloud.EndpointOpts{
		Region:       c.DetermineRegion(region),
		Availability: clientconfig.GetEndpointType(c.EndpointType),
	}
	client, err := newClient(c.OsClient, eo)
	if err, ok := err.(*gophercloud.ErrEndpointNotFound); ok && client != nil {
		client, e := c.DetermineEndpoint(client, eo, service)
		if e != nil {
			return client, e
		}
		// if the endpoint is still not found, return the original error
		if client.ProviderClient == nil {
			return client, err
		}
	}
	if err != nil {
		return client, err
	}

	return c.DetermineEndpoint(client, eo, service)
}

func (c *Config) BlockStorageV1Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewBlockStorageV1, region, "volume")
}

func (c *Config) BlockStorageV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewBlockStorageV2, region, "volumev2")
}

func (c *Config) BlockStorageV3Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewBlockStorageV3, region, "volumev3")
}

func (c *Config) ComputeV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewComputeV2, region, "compute")
}

func (c *Config) DNSV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewDNSV2, region, "dns")
}

func (c *Config) IdentityV3Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewIdentityV3, region, "identity")
}

func (c *Config) ImageV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewImageV2, region, "image")
}

func (c *Config) MessagingV2Client(ctx context.Context, clientID string, region string) (*gophercloud.ServiceClient, error) {
	if err := c.Authenticate(ctx); err != nil {
		return nil, err
	}

	eo := gophercloud.EndpointOpts{
		Region:       c.DetermineRegion(region),
		Availability: clientconfig.GetEndpointType(c.EndpointType),
	}
	client, err := openstack.NewMessagingV2(c.OsClient, clientID, eo)
	if err, ok := err.(*gophercloud.ErrEndpointNotFound); ok && client != nil {
		client, e := c.DetermineEndpoint(client, eo, "messaging")
		if e != nil {
			return client, e
		}
		// if the endpoint is still not found, return the original error
		if client.ProviderClient == nil {
			return client, err
		}
	}
	if err != nil {
		return client, err
	}

	return c.DetermineEndpoint(client, eo, "messaging")
}

func (c *Config) NetworkingV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewNetworkV2, region, "network")
}

func (c *Config) ObjectStorageV1Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	if !c.Swauth {
		return c.CommonServiceClientInit(ctx, openstack.NewObjectStorageV1, region, "object-store")
	}

	// If Swift Authentication is being used, return a swauth client.
	if !c.DelayedAuth {
		return swauth.NewObjectStorageV1(ctx, c.OsClient, swauth.AuthOpts{
			User: c.Username,
			Key:  c.Password,
		})
	}

	c.MutexKV.Lock("SwAuth")
	defer c.MutexKV.Unlock("SwAuth")

	if c.swAuthFailed != nil {
		return nil, c.swAuthFailed
	}

	if c.swClient == nil {
		c.swClient, c.swAuthFailed = swauth.NewObjectStorageV1(ctx, c.OsClient, swauth.AuthOpts{
			User: c.Username,
			Key:  c.Password,
		})
		return c.swClient, c.swAuthFailed
	}

	return c.swClient, nil
}

func (c *Config) OrchestrationV1Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewOrchestrationV1, region, "orchestration")
}

func (c *Config) LoadBalancerV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewLoadBalancerV2, region, "octavia")
}

func (c *Config) DatabaseV1Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewDBV1, region, "database")
}

func (c *Config) ContainerInfraV1Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewContainerInfraV1, region, "container-infra")
}

func (c *Config) SharedfilesystemV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewSharedFileSystemV2, region, "sharev2")
}

func (c *Config) KeyManagerV1Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewKeyManagerV1, region, "key-manager")
}

func (c *Config) WorkflowV2Client(ctx context.Context, region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(ctx, openstack.NewWorkflowV2, region, "workflowv2")
}

// A wrapper to determine if logging in gophercloud should be enabled, with a fallback
// to the OS_DEBUG environment variable when no explicit configuration is passed.
func enableLogging(enable bool) bool {
	if enable {
		return true
	}

	// if OS_DEBUG is set, log the requests and responses
	if os.Getenv("OS_DEBUG") != "" {
		return true
	}

	return false
}
