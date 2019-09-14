package cloudca

import (
	"fmt"
	"testing"

	"github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLoadBalancerRuleCreate(t *testing.T) {
	t.Parallel()

	environmentID := "a225a598-f440-439e-a51e-1c5275bc6d57"
	vpcID := "438fe7a0-d7a6-44f8-875d-b976021a6ae4"
	networkID := "1d5c1e64-59f1-4a34-8539-77af5153058c"
	instanceName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancerRuleCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerRuleCreate(environmentID, vpcID, networkID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerRuleCreateExists("cloudca_load_balancer_rule.foobar"),
				),
			},
			{
				ResourceName:      "cloudca_load_balancer_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoadBalancerRuleCreate(environment, vpc, network, name string) string {
	return fmt.Sprintf(`
resource "cloudca_instance" "foobar" {
	environment_id   = "%s"
	network_id       = "%s"
	name             = "%s"
	template         = "Ubuntu 18.04.2"
	compute_offering = "Standard"
	cpu_count        = 1
	memory_in_mb     = 1024
}

resource "cloudca_public_ip" "foobar" {
	environment_id = "%s"
	vpc_id         = "%s"
}

resource "cloudca_load_balancer_rule" "foobar" {
	environment_id = "%s"
    network_id     = "%s"
    name           = "%s"
	public_ip_id   = "${cloudca_public_ip.foobar.id}"
    protocol       = "tcp"
    algorithm      = "leastconn"
    public_port    = 80
    private_port   = 80
    instance_ids   = ["${cloudca_instance.foobar.id}"]
}`, environment, network, name, environment, vpc, environment, network, name)
}

func testAccCheckLoadBalancerRuleCreateExists(n string) resource.TestCheckFunc {
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

		found, err := resources.LoadBalancerRules.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Load Balancer Rule not found")
		}

		return nil
	}
}

func testAccCheckLoadBalancerRuleCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_load_balancer_rule" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.LoadBalancerRules.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Load Balancer Rule still exists")
			}
		}
	}

	return nil
}
