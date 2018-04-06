package testing

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/openstack/clientconfig"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGetCloudFromYAML(t *testing.T) {
	clientOpts := &clientconfig.ClientOpts{
		Cloud:     "hawaii",
		EnvPrefix: "FOO",
	}

	actual, err := clientconfig.GetCloudFromYAML(clientOpts)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &HawaiiCloudYAML, actual)

	clientOpts = &clientconfig.ClientOpts{
		Cloud:     "california",
		EnvPrefix: "FOO",
	}

	actual, err = clientconfig.GetCloudFromYAML(clientOpts)
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

	clientOpts := &clientconfig.ClientOpts{
		EnvPrefix: "FOO_",
	}

	actual, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertDeepEquals(t, HawaiiAuthOpts, actual)

	os.Unsetenv("FOO_CLOUD")
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
	}

	allClientOpts := map[string]*clientconfig.ClientOpts{
		"hawaii":     HawaiiClientOpts,
		"florida":    FloridaClientOpts,
		"california": CaliforniaClientOpts,
		"arizona":    ArizonaClientOpts,
		"newmexico":  NewMexicoClientOpts,
		"nevada":     NevadaClientOpts,
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
	}

	expectedAuthOpts := map[string]*gophercloud.AuthOptions{
		"hawaii":     HawaiiAuthOpts,
		"florida":    FloridaAuthOpts,
		"california": CaliforniaAuthOpts,
		"arizona":    ArizonaAuthOpts,
		"newmexico":  NewMexicoAuthOpts,
		"nevada":     NevadaAuthOpts,
	}

	for cloud, envVars := range allEnvVars {
		for k, v := range envVars {
			os.Setenv(k, v)
		}

		actualAuthOpts, err := clientconfig.AuthOptions(nil)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expectedAuthOpts[cloud], actualAuthOpts)

		for k, _ := range envVars {
			os.Unsetenv(k)
		}
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
		for k, v := range envVars {
			os.Setenv(k, v)
		}

		actualAuthOpts, err := clientconfig.AuthOptions(nil)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, expectedAuthOpts[cloud], actualAuthOpts)

		for k, _ := range envVars {
			os.Unsetenv(k)
		}
	}
}
