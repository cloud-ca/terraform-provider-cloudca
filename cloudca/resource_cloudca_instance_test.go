package cloudca

import (
	"testing"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

type TestApiClient struct {

}

func (t TestApiClient) Do(request api.CcaRequest) (*api.CcaResponse, error) {
	return nil, nil
}

func testCcaClient() *gocca.CcaClient {
	return gocca.NewCcaClientWithApiClient(TestApiClient{})
}

func testProvider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: GetCloudCAResourceMap(),
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return testCcaClient(), nil
}

func TestUpdateCloudcaInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
			Providers: map[string]terraform.ResourceProvider{
				"cloudca": testProvider(),
			},
		})
}