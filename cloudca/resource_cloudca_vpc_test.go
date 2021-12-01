package cloudca

import (
	"fmt"
	"testing"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVPCCreate(t *testing.T) {
	t.Parallel()

	vpcName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCCreate(environmentID, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCCreateExists("cloudca_vpc.foobar"),
				),
			},
		},
	})
}

func testAccVPCCreate(environment, name string) string {
	return fmt.Sprintf(`
resource "cloudca_vpc" "foobar" {
	environment_id = "%s"
	name           = "%s"
	description    = "This is a %s vpc"
	vpc_offering   = "Default VPC offering"
}`, environment, name, name)
}

func testAccCheckVPCCreateExists(n string) resource.TestCheckFunc {
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

		found, err := resources.Vpcs.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("VPC not found")
		}

		return nil
	}
}

func testAccCheckVPCCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_vpc" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.Vpcs.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("VPC still exists")
			}
		}
	}

	return nil
}
