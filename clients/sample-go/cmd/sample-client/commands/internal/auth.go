package internal

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/credentials"
	managerapi "github.com/uor-framework/uor-client-go/api/services/collectionmanager/v1alpha1"
)

// GetCredentials returns an AuthConfig resource for a given reference.
func GetCredentials(reference string) (*managerapi.AuthConfig, error) {
	uri, err := url.Parse("https://" + reference)
	if err != nil {
		return nil, err
	}
	cfg, err := loadDefaultConfig()
	if err != nil {
		return nil, err
	}

	authConf, err := cfg.GetCredentialsStore(uri.Host).Get(uri.Host)
	if err != nil {
		return nil, err
	}
	cred := managerapi.AuthConfig{
		ServerAddress: uri.String(),
		Username:      authConf.Username,
		Password:      authConf.Password,
		AccessToken:   authConf.RegistryToken,
		RefreshToken:  authConf.IdentityToken,
	}
	return &cred, nil
}

// loadConfigFile reads the credential-related configuration
// from the given path.
func loadConfigFile(path string) (*configfile.ConfigFile, error) {
	cfg := configfile.New(path)
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return cfg, err
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := cfg.LoadFromReader(file); err != nil {
		return nil, err
	}

	if !cfg.ContainsAuth() {
		cfg.CredentialsStore = credentials.DetectDefaultStore(cfg.CredentialsStore)
	}

	return cfg, nil
}

// LoadDefaultConfig attempts to load credentials in the
// default Docker location, then the default Podman location.
func loadDefaultConfig() (*configfile.ConfigFile, error) {
	dir := config.Dir()
	dockerConfigJSON := filepath.Join(dir, config.ConfigFileName)
	cfg := configfile.New(dockerConfigJSON)

	switch _, err := os.Stat(dockerConfigJSON); {
	case err == nil:
		cfg, err = config.Load(dir)
		if err != nil {
			return cfg, err
		}
	case os.IsNotExist(err):
		podmanConfig := filepath.Join(xdg.RuntimeDir, "containers/auth.json")
		cfg, err = loadConfigFile(podmanConfig)
		if err != nil {
			return cfg, err
		}
	}

	if !cfg.ContainsAuth() {
		cfg.CredentialsStore = credentials.DetectDefaultStore(cfg.CredentialsStore)
	}

	return cfg, nil
}
