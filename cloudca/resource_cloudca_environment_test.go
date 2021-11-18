package cloudca

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/cloud-ca/go-cloudca"
)

func TestAccEnvironmentCreate(t *testing.T) {
	t.Parallel()

	environmentName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentCreate(environmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentCreateExists("cloudca_environment.foobar"),
				),
			},
			{
				ResourceName:      "cloudca_environment.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEnvironmentCreate(name string) string {
	return fmt.Sprintf(`
resource "cloudca_environment" "foobar" {
	organization_code = "lab"
	service_code      = "compute-on"
	name              = "%s"
	description       = "Environment for %s workloads"
	admin_role        = []
	read_only_role    = []
}`, name, name)
}

func testAccCheckEnvironmentCreateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*cca.CcaClient)

		found, err := client.Environments.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Environment not found")
		}

		return nil
	}
}

func testAccCheckEnvironmentCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_environment" {
			_, err := client.Environments.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Environment still exists")
			}
		}
	}

	return nil
}
