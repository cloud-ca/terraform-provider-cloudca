package cloudca

import (
	"fmt"
	"testing"

	"github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPublicIPCreate(t *testing.T) {
	t.Parallel()

	environmentID := "a225a598-f440-439e-a51e-1c5275bc6d57"
	vpcID := "438fe7a0-d7a6-44f8-875d-b976021a6ae4"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPublicIPCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIPCreate(environmentID, vpcID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicIPCreateExists("cloudca_public_ip.foobar"),
				),
			},
			{
				ResourceName:      "cloudca_public_ip.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPublicIPCreate(environment, vpc string) string {
	return fmt.Sprintf(`
resource "cloudca_public_ip" "foobar" {
	environment_id = "%s"
	vpc_id         = "%s"
}`, environment, vpc)
}

func testAccCheckPublicIPCreateExists(n string) resource.TestCheckFunc {
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

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Public IP not found")
		}

		return nil
	}
}

func testAccCheckPublicIPCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_public_ip" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.PublicIps.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Public IP still exists")
			}
		}
	}

	return nil
}
