package main

import "github.com/cloud-ca/go-cloudca" 

// Config is the configuration structure used to instantiate a
// new CloudStack client.
type Config struct {
	APIURL string
	APIKey string
	Insecure bool
}

// NewClient returns a new CcaClient client.
func (c *Config) NewClient() (*gocca.CcaClient, error) {
	if c.Insecure {
		return gocca.NewInsecureCcaClientWithURL(c.APIURL, c.APIKey), nil
	}
	return gocca.NewCcaClientWithURL(c.APIURL, c.APIKey), nil
}
