package testing

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	yaml "gopkg.in/yaml.v3"
)

func TestGetCloudFromYAML(t *testing.T) {
	allClientOpts := map[string]*clientconfig.ClientOpts{
		"hawaii": {
			Cloud:     "hawaii",
			EnvPrefix: "FOO",
		},
		"california": {
			Cloud:     "california",
			EnvPrefix: "FOO",
		},
		"florida_insecure":   {Cloud: "florida_insecure"},
		"florida_secure":     {Cloud: "florida_secure"},
		"nevada":             {Cloud: "nevada"},
		"texas":              {Cloud: "texas"},
		"alberta":            {Cloud: "alberta"},
		"yukon":              {Cloud: "yukon"},
		"chicago":            {Cloud: "chicago"},
		"chicago_legacy":     {Cloud: "chicago_legacy"},
		"chicago_useprofile": {Cloud: "chicago_useprofile"},
		"philadelphia":       {Cloud: "philadelphia"},
		"philadelphia_phl1": {
			Cloud:      "philadelphia_complex",
			RegionName: "PHL1",
		},
		"philadelphia_phl2": {
			Cloud:      "philadelphia_complex",
			RegionName: "PHL2",
		},
		"virginia": {Cloud: "virginia"},
		"disconnected_smw": {
			Cloud:      "disconnected_clouds",
			RegionName: "SOMEWHERE",
		},
		"disconnected_anw": {
			Cloud:      "disconnected_clouds",
			RegionName: "ANYWHERE",
		},
		"disconnected_now": {
			Cloud:      "disconnected_clouds",
			RegionName: "NOWHERE",
		},
	}

	expectedClouds := map[string]*clientconfig.Cloud{
		"hawaii":             &HawaiiCloudYAML,
		"california":         &CaliforniaCloudYAML,
		"florida_insecure":   &InsecureFloridaCloudYAML,
		"florida_secure":     &SecureFloridaCloudYAML,
		"nevada":             &NevadaCloudYAML,
		"texas":              &TexasCloudYAML,
		"alberta":            &AlbertaCloudYAML,
		"yukon":              &YukonCloudYAML,
		"chicago":            &ChicagoCloudYAML,
		"chicago_legacy":     &ChicagoCloudLegacyYAML,
		"chicago_useprofile": &ChicagoCloudUseProfileYAML,
		"philadelphia":       &PhiladelphiaCloudYAML,
		"philadelphia_phl1":  &PhiladelphiaComplexPhl1CloudYAML,
		"philadelphia_phl2":  &PhiladelphiaComplexPhl2CloudYAML,
		"virginia":           &VirginiaCloudYAML,
		"disconnected_smw":   &DisconnectedSomewhereCloudYAML,
		"disconnected_anw":   &DisconnectedAnywhereCloudYAML,
		"disconnected_now":   &DisconnectedNowhereCloudYAML,
	}

	for cloud, clientOpts := range allClientOpts {
		actual, err := clientconfig.GetCloudFromYAML(clientOpts)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expectedClouds[cloud], actual)
	}
}

func TestGetCloudFromYAMLOSCLOUD(t *testing.T) {
	os.Setenv("OS_CLOUD", "california")
	defer os.Unsetenv("OS_CLOUD")

	clientOpts := &clientconfig.ClientOpts{
		Cloud: "hawaii",
	}

	actual, err := clientconfig.GetCloudFromYAML(clientOpts)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &HawaiiCloudYAML, actual)
}

func TestGetCloudFromYAMLMissingClientOpts(t *testing.T) {
	os.Setenv("OS_CLOUD", "california")
	defer os.Unsetenv("OS_CLOUD")

	actual, err := clientconfig.GetCloudFromYAML(nil)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &CaliforniaCloudYAML, actual)
}

func TestAuthOptionsExplicitCloud(t *testing.T) {
	os.Unsetenv("OS_CLOUD")

	clientOpts := &clientconfig.ClientOpts{
		Cloud: "hawaii",
	}

	actual, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertDeepEquals(t, HawaiiAuthOpts, actual)
}

func TestAuthOptionsOSCLOUD(t *testing.T) {
	os.Setenv("FOO_CLOUD", "hawaii")
	defer os.Unsetenv("FOO_CLOUD")

	clientOpts := &clientconfig.ClientOpts{
		EnvPrefix: "FOO_",
	}

	actual, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertDeepEquals(t, HawaiiAuthOpts, actual)
}

func TestAuthOptionsExplicitCloudAndOSCLOUD(t *testing.T) {
	os.Setenv("FOO_CLOUD", "hawaii")
	defer os.Unsetenv("FOO_CLOUD")

	clientOpts := &clientconfig.ClientOpts{
		EnvPrefix: "FOO_",
		Cloud:     "california",
	}

	actual, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		t.Fatal(err)
	}

	// We should have ignored the cloud configuration option
	th.AssertDeepEquals(t, CaliforniaAuthOpts, actual)
}

func TestAuthOptionsMissingClientOpts(t *testing.T) {
	os.Setenv("OS_CLOUD", "hawaii")
	defer os.Unsetenv("OS_CLOUD")

	actual, err := clientconfig.AuthOptions(nil)
	if err != nil {
		t.Fatal(err)
	}

	// We should have handled the missing config opts and fallen back to
	// defaults
	th.AssertDeepEquals(t, HawaiiAuthOpts, actual)
}

func TestAuthOptionsCreationFromCloudsYAML(t *testing.T) {
	os.Unsetenv("OS_CLOUD")

	allClouds := map[string]*gophercloud.AuthOptions{
		"hawaii":     HawaiiAuthOpts,
		"florida":    FloridaAuthOpts,
		"california": CaliforniaAuthOpts,
		"arizona":    ArizonaAuthOpts,
		"newmexico":  NewMexicoAuthOpts,
		"nevada":     NevadaAuthOpts,
		"texas":      TexasAuthOpts,
	}

	for cloud, expected := range allClouds {
		clientOpts := &clientconfig.ClientOpts{
			Cloud: cloud,
		}

		actual, err := clientconfig.AuthOptions(clientOpts)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expected, actual)

		scope, err := expected.ToTokenV3ScopeMap()
		th.AssertNoErr(t, err)

		_, err = expected.ToTokenV3CreateMap(scope)
		th.AssertNoErr(t, err)
	}
}

func TestAuthOptionsCreationFromLegacyCloudsYAML(t *testing.T) {
	os.Unsetenv("OS_CLOUD")

	allClouds := map[string]*gophercloud.AuthOptions{
		"alberta": AlbertaAuthOpts,
		"yukon":   YukonAuthOpts,
	}

	for cloud, expected := range allClouds {
		clientOpts := &clientconfig.ClientOpts{
			Cloud: cloud,
		}

		actual, err := clientconfig.AuthOptions(clientOpts)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expected, actual)

		_, err = expected.ToTokenV2CreateMap()
		th.AssertNoErr(t, err)
	}
}

func TestAuthOptionsCreationFromClientConfig(t *testing.T) {
	os.Unsetenv("OS_CLOUD")

	expectedAuthOpts := map[string]*gophercloud.AuthOptions{
		"hawaii":     HawaiiAuthOpts,
		"florida":    FloridaAuthOpts,
		"california": CaliforniaAuthOpts,
		"arizona":    ArizonaAuthOpts,
		"newmexico":  NewMexicoAuthOpts,
		"nevada":     NevadaAuthOpts,
		"texas":      TexasAuthOpts,
	}

	allClientOpts := map[string]*clientconfig.ClientOpts{
		"hawaii":     HawaiiClientOpts,
		"florida":    FloridaClientOpts,
		"california": CaliforniaClientOpts,
		"arizona":    ArizonaClientOpts,
		"newmexico":  NewMexicoClientOpts,
		"nevada":     NevadaClientOpts,
		"texas":      TexasClientOpts,
	}

	for cloud, clientOpts := range allClientOpts {
		actualAuthOpts, err := clientconfig.AuthOptions(clientOpts)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expectedAuthOpts[cloud], actualAuthOpts)
	}
}

func TestAuthOptionsCreationFromLegacyClientConfig(t *testing.T) {
	os.Unsetenv("OS_CLOUD")

	expectedAuthOpts := map[string]*gophercloud.AuthOptions{
		"alberta": AlbertaAuthOpts,
		"yukon":   YukonAuthOpts,
	}

	allClientOpts := map[string]*clientconfig.ClientOpts{
		"alberta": AlbertaClientOpts,
		"yukon":   YukonClientOpts,
	}

	for cloud, clientOpts := range allClientOpts {
		actualAuthOpts, err := clientconfig.AuthOptions(clientOpts)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expectedAuthOpts[cloud], actualAuthOpts)
	}
}

func TestAuthOptionsCreationFromEnv(t *testing.T) {
	os.Unsetenv("OS_CLOUD")

	allEnvVars := map[string]map[string]string{
		"hawaii":     HawaiiEnvAuth,
		"florida":    FloridaEnvAuth,
		"california": CaliforniaEnvAuth,
		"arizona":    ArizonaEnvAuth,
		"newmexico":  NewMexicoEnvAuth,
		"nevada":     NevadaEnvAuth,
		"texas":      TexasEnvAuth,
		"virginia":   VirginiaEnvAuth,
	}

	expectedAuthOpts := map[string]*gophercloud.AuthOptions{
		"hawaii":     HawaiiAuthOpts,
		"florida":    FloridaAuthOpts,
		"california": CaliforniaAuthOpts,
		"arizona":    ArizonaAuthOpts,
		"newmexico":  NewMexicoAuthOpts,
		"nevada":     NevadaAuthOpts,
		"texas":      TexasAuthOpts,
		"virginia":   VirginiaAuthOpts,
	}

	for cloud, envVars := range allEnvVars {
		t.Run(cloud, func(t *testing.T) {
			for k, v := range envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			actualAuthOpts, err := clientconfig.AuthOptions(nil)
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, expectedAuthOpts[cloud], actualAuthOpts)
		})
	}
}

func TestAuthOptionsCreationFromLegacyEnv(t *testing.T) {
	os.Unsetenv("OS_CLOUD")

	allEnvVars := map[string]map[string]string{
		"alberta": AlbertaEnvAuth,
		"yukon":   YukonEnvAuth,
	}

	expectedAuthOpts := map[string]*gophercloud.AuthOptions{
		"alberta": AlbertaAuthOpts,
		"yukon":   YukonAuthOpts,
	}

	for cloud, envVars := range allEnvVars {
		t.Run(cloud, func(t *testing.T) {
			for k, v := range envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			actualAuthOpts, err := clientconfig.AuthOptions(nil)
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, expectedAuthOpts[cloud], actualAuthOpts)
		})
	}
}

type CustomYAMLOpts struct {
	Logger *log.Logger
}

func (opts CustomYAMLOpts) LoadCloudsYAML() (map[string]clientconfig.Cloud, error) {
	filename, content, err := clientconfig.FindAndReadCloudsYAML()
	if err != nil {
		return nil, err
	}

	var clouds clientconfig.Clouds
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	opts.Logger.Printf("Filename: %s", filename)

	return clouds.Clouds, nil
}

func (opts CustomYAMLOpts) LoadSecureCloudsYAML() (map[string]clientconfig.Cloud, error) {
	return nil, nil
}

func (opts CustomYAMLOpts) LoadPublicCloudsYAML() (map[string]clientconfig.Cloud, error) {
	return nil, nil
}

func TestGetCloudFromYAMLWithCustomYAMLOpts(t *testing.T) {
	logger := log.New(os.Stderr, "", log.Lshortfile)
	yamlOpts := CustomYAMLOpts{
		Logger: logger,
	}

	allClientOpts := map[string]*clientconfig.ClientOpts{
		"hawaii": {
			Cloud:     "hawaii",
			EnvPrefix: "FOO",
			YAMLOpts:  yamlOpts,
		},
		"california": {
			Cloud:     "california",
			EnvPrefix: "FOO",
			YAMLOpts:  yamlOpts,
		},
	}

	expectedClouds := map[string]*clientconfig.Cloud{
		"hawaii":     &HawaiiCloudYAML,
		"california": &CaliforniaCloudYAML,
	}

	for cloud, clientOpts := range allClientOpts {
		actual, err := clientconfig.GetCloudFromYAML(clientOpts)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expectedClouds[cloud], actual)
	}
}
