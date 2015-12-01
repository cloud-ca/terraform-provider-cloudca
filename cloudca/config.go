package cloudca

import "github.com/cloud-ca/go-cloudca" 

// Config is the configuration structure used to instantiate a
// new CloudStack client.
type Config struct {
	APIURL string
	APIKey string
}

// NewClient returns a new CcaClient client.
func (c *Config) NewClient() (*gocca.CcaClient, error) {
	cca := gocca.NewCcaClientWithURL(c.APIURL, c.APIKey)
	return cca, nil
}
