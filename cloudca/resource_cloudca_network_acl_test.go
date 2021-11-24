package cloudca

import (
	"fmt"
	"testing"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkACLCreate(t *testing.T) {
	t.Parallel()

	environmentID := "c67a090f-b66f-42e1-b444-10cdff9d8be2"
	vpcID := "2c01d952-d010-4811-b66d-4c7f5f805193"
	networkACLName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLCreate(environmentID, vpcID, networkACLName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLCreateExists("cloudca_network_acl.foobar"),
				),
			},
		},
	})
}

func testAccNetworkACLCreate(environment, vpc, name string) string {
	return fmt.Sprintf(`
resource "cloudca_network_acl" "foobar" {
	environment_id = "%s"
	vpc_id         = "%s"
	name           = "%s"
	description    = "This is a %s acl"
}`, environment, vpc, name, name)
}

func testAccCheckNetworkACLCreateExists(n string) resource.TestCheckFunc {
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

		found, err := resources.NetworkAcls.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Network ACL not found")
		}

		return nil
	}
}

func testAccCheckNetworkACLCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_network_acl" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.NetworkAcls.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Network ACL still exists")
			}
		}
	}

	return nil
}
