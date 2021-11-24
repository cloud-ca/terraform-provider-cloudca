package cloudca

import cca "github.com/cloud-ca/go-cloudca"

// Config is the configuration structure used to instantiate a
// new cloudca client.
type Config struct {
	APIURL   string
	APIKey   string
	Insecure bool
}

// NewClient returns a new CcaClient client.
func (c *Config) NewClient() (*cca.CcaClient, error) {
	if c.Insecure {
		return cca.NewInsecureCcaClientWithURL(c.APIURL, c.APIKey), nil
	}
	return cca.NewCcaClientWithURL(c.APIURL, c.APIKey), nil
}
