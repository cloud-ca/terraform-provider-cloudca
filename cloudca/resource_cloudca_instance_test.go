package cloudca

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/cloud-ca/go-cloudca"
)

func TestAccInstanceCreateBasic(t *testing.T) {
	t.Parallel()

	environmentID := "c67a090f-b66f-42e1-b444-10cdff9d8be2"
	networkID := "719af2c3-2da8-474f-b03e-63fce6e1a827"
	instanceName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceCreateBasicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceCreateBasic(environmentID, networkID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceCreateBasicExists("cloudca_instance.foobar"),
				),
			},
			{
				ResourceName:      "cloudca_instance.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceCreateDataDrive(t *testing.T) {
	t.Parallel()

	environmentID := "c67a090f-b66f-42e1-b444-10cdff9d8be2"
	networkID := "719af2c3-2da8-474f-b03e-63fce6e1a827"
	instanceName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceCreateDataDriveDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceCreateDataDrive(environmentID, networkID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceCreateDataDriveExists("cloudca_instance.foobar"),
				),
			},
			{
				ResourceName:      "cloudca_instance.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInstanceCreateBasic(environment, network, name string) string {
	return fmt.Sprintf(`
resource "cloudca_instance" "foobar" {
	environment_id   = "%s"
	network_id       = "%s"
	name             = "%s"
	template         = "Ubuntu 20.04.2"
	compute_offering = "Standard"
	cpu_count        = 1
	memory_in_mb     = 1024
}`, environment, network, name)
}

func testAccInstanceCreateDataDrive(environment, network, name string) string {
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
resource "cloudca_volume" "foobar" {
	environment_id = "%s"
    name           = "%s"
	disk_offering  = "Performance, No QoS"
	size_in_gb     = "10"
    instance_id    = "${cloudca_instance.foobar.id}"
}`, environment, network, name, environment, name)
}

func testAccCheckInstanceCreateBasicExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
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

		found, err := resources.Instances.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID || found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Instance not found")
		}

		return nil
	}
}

func testAccCheckInstanceCreateDataDriveExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
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

		found, err := resources.Instances.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID || found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Instance not found")
		}

		return nil
	}
}

func testAccCheckInstanceCreateBasicDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_instance" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.Instances.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckInstanceCreateDataDriveDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_instance" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.Instances.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}
