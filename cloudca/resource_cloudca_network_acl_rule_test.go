package cloudca

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/cloud-ca/go-cloudca"
)

func TestAccNetworkACLRuleCreate(t *testing.T) {
	t.Parallel()

	environmentID := "c67a090f-b66f-42e1-b444-10cdff9d8be2"
	vpcID := "2c01d952-d010-4811-b66d-4c7f5f805193" 
	networkACLRuleName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLRuleCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRuleCreate(environmentID, vpcID, networkACLRuleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleCreateExists("cloudca_network_acl_rule.foobar"),
				),
			},
			{
				ResourceName:      "cloudca_network_acl_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkACLRuleCreate(environment, vpc, name string) string {
	return fmt.Sprintf(`
resource "cloudca_network_acl" "foobar" {
	environment_id = "%s"
	vpc_id         = "%s"
	name           = "%s"
	description    = "This is a %s acl"
}
resource "cloudca_network_acl_rule" "foobar" {
	environment_id = "%s"
	network_acl_id = "${cloudca_network_acl.foobar.id}"
	rule_number    = 55
	cidr           = "10.212.208.0/22"
	action         = "Allow"
	protocol       = "TCP"
	start_port     = 80
	end_port       = 80
	traffic_type   = "Ingress"
}`, environment, vpc, name, name, environment)
}

func testAccCheckNetworkACLRuleCreateExists(n string) resource.TestCheckFunc {
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

		found, err := resources.NetworkAclRules.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Network ACL Rule not found")
		}

		return nil
	}
}

func testAccCheckNetworkACLRuleCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_network_acl_rule" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.NetworkAclRules.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Network ACL Rule still exists")
			}
		}
	}

	return nil
}
