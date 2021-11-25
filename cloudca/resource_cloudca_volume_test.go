package cloudca

import (
	"fmt"
	"testing"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolumeCreate(t *testing.T) {
	t.Parallel()

	instanceID := "6f26111d-464d-4fc8-9c72-7a181a96c257"
	volumeName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeCreate(environmentID, instanceID, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeCreateExists("cloudca_volume.foobar"),
				),
			},
		},
	})
}

func testAccVolumeCreate(environment, instance, name string) string {
	return fmt.Sprintf(`
resource "cloudca_volume" "foobar" {
	environment_id = "%s"
	name           = "%s"
	disk_offering  = "fd78763c-f33a-43f3-b1e3-63bf59a48350"
    instance_id    = "%s"
	size_in_gb     = "10"
}`, environment, name, instance)
}

func testAccCheckVolumeCreateExists(n string) resource.TestCheckFunc {
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

		found, err := resources.Volumes.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		return nil
	}
}

func testAccCheckVolumeCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_volume" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.Volumes.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Volume still exists")
			}
		}
	}

	return nil
}
