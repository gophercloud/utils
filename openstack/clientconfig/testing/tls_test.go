package testing

import (
	"os"
	"testing"

	"github.com/gophercloud/utils/v2/openstack/clientconfig"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestPrepareTLSConfigInsecureEnv(t *testing.T) {
	t.Run("OS_INSECURE=true", func(t *testing.T) {
		os.Setenv("OS_INSECURE", "true")
		defer os.Unsetenv("OS_INSECURE")

		tlsConfig, err := clientconfig.PrepareTLSConfig("OS_", &clientconfig.Cloud{})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, true, tlsConfig.InsecureSkipVerify)
	})

	t.Run("OS_INSECURE=false", func(t *testing.T) {
		os.Setenv("OS_INSECURE", "false")
		defer os.Unsetenv("OS_INSECURE")

		tlsConfig, err := clientconfig.PrepareTLSConfig("OS_", &clientconfig.Cloud{})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, false, tlsConfig.InsecureSkipVerify)
	})

	t.Run("OS_INSECURE unset", func(t *testing.T) {
		os.Unsetenv("OS_INSECURE")

		tlsConfig, err := clientconfig.PrepareTLSConfig("OS_", &clientconfig.Cloud{})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, false, tlsConfig.InsecureSkipVerify)
	})

	t.Run("OS_INSECURE=invalid", func(t *testing.T) {
		os.Setenv("OS_INSECURE", "invalid")
		defer os.Unsetenv("OS_INSECURE")

		_, err := clientconfig.PrepareTLSConfig("OS_", &clientconfig.Cloud{})
		th.AssertErr(t, err)
	})

	t.Run("cloud.Verify overrides OS_INSECURE", func(t *testing.T) {
		os.Setenv("OS_INSECURE", "true")
		defer os.Unsetenv("OS_INSECURE")

		cloud := &clientconfig.Cloud{Verify: &iTrue}
		tlsConfig, err := clientconfig.PrepareTLSConfig("OS_", cloud)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, false, tlsConfig.InsecureSkipVerify)
	})

	t.Run("custom env prefix", func(t *testing.T) {
		os.Setenv("FOO_INSECURE", "true")
		defer os.Unsetenv("FOO_INSECURE")

		tlsConfig, err := clientconfig.PrepareTLSConfig("FOO_", &clientconfig.Cloud{})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, true, tlsConfig.InsecureSkipVerify)
	})
}
