package cloudca

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cloudcaAPIKey = "CLOUDCA_API_KEY"

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cloudca": testAccProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func hasEnvValue(envKey string) bool {
	if v := os.Getenv(envKey); v == "" {
		return false
	}
	return true
}

func testAccPreCheckEnvs(t *testing.T, keys ...string) {
	for _, env := range keys {
		if !hasEnvValue(env) {
			t.Fatal(fmt.Sprintf("%s must be set for acceptance tests", env))
		}
	}

}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckEnvs(t, cloudcaAPIKey)
}
