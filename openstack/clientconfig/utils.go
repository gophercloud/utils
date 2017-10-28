package clientconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Wessie/appdirs"
)

// findAndLoadYAML attempts to locate a clouds.yaml file in the following
// locations:
//
// 1. OS_CLIENT_CONFIG_FILE
// 2. Current directory.
// 3. XDG OS-specific user_config_dir
// 4. unix-specific user_config_dir (~/.config/openstack/clouds.yaml)
// 5. XDG OS-specific site_config_dir
// 6. unix-specific site_config_dir (/etc/openstack/clouds.yaml)
//
// If found, the contents of the file is returned.
func findAndReadYAML() ([]byte, error) {
	// OS_CLIENT_CONFIG_FILE
	if v := os.Getenv("OS_CLIENT_CONFIG_FILE"); v != "" {
		if ok := fileExists(v); ok {
			return ioutil.ReadFile(v)
		}
	}

	// current directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to determine working directory: %s", err)
	}

	filename := filepath.Join(cwd, "clouds.yaml")
	if ok := fileExists(filename); ok {
		return ioutil.ReadFile(filename)
	}

	// xdg os-specific user_config_dir
	app := appdirs.New("openstack", "", "")
	if v := app.UserConfig(); v != "" {
		filename := filepath.Join(v, "clouds.yaml")
		if ok := fileExists(filename); ok {
			return ioutil.ReadFile(filename)
		}
	}

	// unix-specific user_config_dir.
	// xdg on Mac/Darwin is "Application Support",
	// but maybe ~/.config/openstack exists anyway.
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to get current user: %s", err)
	}

	homeDir := currentUser.HomeDir
	if homeDir != "" {
		filename := filepath.Join(homeDir, ".config/openstack/clouds.yaml")
		if ok := fileExists(filename); ok {
			return ioutil.ReadFile(filename)
		}
	}

	// xdg OS-specific site_config_dir
	if v := app.SiteConfig(); v != "" {
		filename := filepath.Join(v, "clouds.yaml")
		if ok := fileExists(filename); ok {
			return ioutil.ReadFile(filename)
		}
	}

	// unix-specific site_config_dir
	if ok := fileExists("/etc/openstack/clouds.yaml"); ok {
		return ioutil.ReadFile("/etc/openstack/clouds.yaml")
	}

	return nil, fmt.Errorf("no clouds.yaml file found")
}

// fileExists checks for the existence of a file at a given location.
func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
