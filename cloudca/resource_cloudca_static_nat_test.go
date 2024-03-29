package cloudca

import (
	"fmt"
	"testing"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStaticNATCreate(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStaticNATCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticNATCreate(environmentID, vpcID, networkID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStaticNATCreateExists("cloudca_static_nat.foobar"),
				),
			},
		},
	})
}

func testAccStaticNATCreate(environment, vpc, network, name string) string {
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
resource "cloudca_static_nat" "foobar" {
	environment_id = "%s"
	public_ip_id   = "${cloudca_public_ip.foobar.id}"
	private_ip_id  = "${cloudca_instance.foobar.private_ip_id}"
}`, environment, network, name, environment, vpc, environment)
}

func testAccCheckStaticNATCreateExists(n string) resource.TestCheckFunc {
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

		found, err := resources.PublicIps.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID || found.PrivateIpId == "" {
			return fmt.Errorf("Static NAT not found")
		}

		return nil
	}
}

func testAccCheckStaticNATCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_static_nat" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.PublicIps.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Static NAT still exists")
			}
		}
	}

	return nil
}
