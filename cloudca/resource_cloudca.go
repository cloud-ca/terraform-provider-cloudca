package cloudca

import (
	"fmt"
	"log"
	"regexp"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetCloudCAResourceMap return the available Resource map
func GetCloudCAResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"cloudca_instance":             resourceCloudcaInstance(),
		"cloudca_environment":          resourceCloudcaEnvironment(),
		"cloudca_vpc":                  resourceCloudcaVpc(),
		"cloudca_network":              resourceCloudcaNetwork(),
		"cloudca_port_forwarding_rule": resourceCloudcaPortForwardingRule(),
		"cloudca_public_ip":            resourceCloudcaPublicIP(),
		"cloudca_volume":               resourceCloudcaVolume(),
		"cloudca_load_balancer_rule":   resourceCloudcaLoadBalancerRule(),
		"cloudca_network_acl":          resourceCloudcaNetworkACL(),
		"cloudca_network_acl_rule":     resourceCloudcaNetworkACLRule(),
		"cloudca_static_nat":           resourceCloudcaStaticNAT(),
		"cloudca_ssh_key":              resourceCloudcaSSHKey(),
		"cloudca_vpn":                  resourceCloudcaVpn(),
		"cloudca_vpn_user":             resourceCloudcaVpnUser(),
	}
}

func setValueOrID(d *schema.ResourceData, key string, value string, id string) error {
	if isID(d.Get(key).(string)) {
		return d.Set(key, id)
	}
	return d.Set(key, value)
}

func isID(id string) bool {
	re := regexp.MustCompile(`^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)
	return re.MatchString(id)
}

// Provides a common, simple way to deal with 404s.
func handleNotFoundError(entity string, deleted bool, err error, d *schema.ResourceData) error {
	if ccaError, ok := err.(api.CcaErrorResponse); ok {
		if ccaError.StatusCode == 404 {
			d.SetId("")
			if deleted {
				log.Printf("%s (id=%s) not found", entity, d.Id())
				return nil
			}
			return fmt.Errorf("%s (id=%s) not found", entity, d.Id())
		}
	}
	return err
}

// Deals with all of the casting done to get a cloudca.Resources.
func getResources(d *schema.ResourceData, meta interface{}) cloudca.Resources {
	client := meta.(*cca.CcaClient)
	_resources, _ := client.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	return _resources.(cloudca.Resources)
}

// Deals with all of the casting done to get a cloudca.Resources.
func getResourcesForEnvironmentID(client *cca.CcaClient, environmentID string) (cloudca.Resources, error) {
	environment, err := client.Environments.Get(environmentID)
	if err != nil {
		return cloudca.Resources{}, err
	}
	resources, err := client.GetResources(environment.ServiceConnection.ServiceCode, environment.Name)
	if err != nil {
		return cloudca.Resources{}, err
	}
	return resources.(cloudca.Resources), nil
}
