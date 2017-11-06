package testing

import (
	"os"
	"testing"

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
	th.AssertDeepEquals(t, &CloudYAMLHawaii, actual)

	clientOpts = &clientconfig.ClientOpts{
		Cloud:     "california",
		EnvPrefix: "FOO",
	}

	actual, err = clientconfig.GetCloudFromYAML(clientOpts)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &CloudYAMLCalifornia, actual)
}

func TestAuthOptionsExplicitCloud(t *testing.T) {
	clientOpts := &clientconfig.ClientOpts{
		Cloud:     "hawaii",
		EnvPrefix: "FOO",
	}

	actual, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertDeepEquals(t, HawaiiAuthOpts, actual)
}

func TestAuthOptionsOSCLOUD(t *testing.T) {
	os.Setenv("OS_CLOUD", "hawaii")

	clientOpts := &clientconfig.ClientOpts{
		EnvPrefix: "FOO",
	}

	actual, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		t.Fatal(err)
	}

	th.AssertDeepEquals(t, HawaiiAuthOpts, actual)
}
