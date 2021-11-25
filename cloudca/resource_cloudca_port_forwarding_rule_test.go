package cloudca

import (
	"fmt"
	"testing"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPortForwardingRuleCreate(t *testing.T) {
	t.Parallel()

	environmentID := "c67a090f-b66f-42e1-b444-10cdff9d8be2"
	vpcID := "2c01d952-d010-4811-b66d-4c7f5f805193"
	networkID := "405e35c3-3e69-4e02-a162-a4112d94acd9"
	instanceName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPortForwardingRuleCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardingRuleCreate(environmentID, vpcID, networkID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortForwardingRuleCreateExists("cloudca_port_forwarding_rule.foobar"),
				),
			},
		},
	})
}

func testAccPortForwardingRuleCreate(environment, vpc, network, name string) string {
	return fmt.Sprintf(`
resource "cloudca_instance" "foobar" {
	environment_id   = "%s"
	network_id       = "%s"
	name             = "%s"
	template         = "Ubuntu 20.04.2"
	compute_offering = "Standard"
	cpu_count        = 1
	memory_in_mb     = 1024
}
resource "cloudca_public_ip" "foobar" {
	environment_id = "%s"
	vpc_id         = "%s"
}
resource "cloudca_port_forwarding_rule" "foobar" {
	environment_id     = "%s"
	public_ip_id       = "${cloudca_public_ip.foobar.id}"
    public_port_start  = 80
	private_ip_id      = "${cloudca_instance.foobar.private_ip_id}"
    private_port_start = 8080
    protocol           = "TCP"
}`, environment, network, name, environment, vpc, environment)
}

func testAccCheckPortForwardingRuleCreateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["environment_id"] == "" {
			return fmt.Errorf("Environment ID is missing")
		}

		client := testAccProvider.Meta().(*cca.CcaClient)
		resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		found, err := resources.PortForwardingRules.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Port Forwarding Rule not found")
		}

		return nil
	}
}

func testAccCheckPortForwardingRuleCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_port_forwarding_rule" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.PortForwardingRules.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Port Forwarding Rule still exists")
			}
		}
	}

	return nil
}
