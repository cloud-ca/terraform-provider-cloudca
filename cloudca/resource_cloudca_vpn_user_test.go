package cloudca

import (
	"fmt"
	"testing"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRemoteAccessVPNUserCreate(t *testing.T) {
	/*
	 this tests conflicts with the VPN Enable test, since it
	 expects the VPN to be disabled,
	*/
	// t.Parallel()

	environmentID := "c67a090f-b66f-42e1-b444-10cdff9d8be2"
	vpcID := "2c01d952-d010-4811-b66d-4c7f5f805193"
	vpnUserName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRemoteAccessVPNUserCreateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRemoteAccessVPNUserCreate(environmentID, vpcID, vpnUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRemoteAccessVPNUserCreateExists("cloudca_vpn_user.foobar"),
				),
			},
		},
	})
}

func testAccRemoteAccessVPNUserCreate(environment, vpc, username string) string {
	return fmt.Sprintf(`
resource "cloudca_vpn" "foobar" {
	environment_id = "%s"
	vpc_id         = "%s"
}
resource "cloudca_vpn_user" "foobar" {
	environment_id = "%s"
	username       = "%s"
	password       = "foopassword"
}`, environment, vpc, environment, username)
}

func testAccCheckRemoteAccessVPNUserCreateExists(n string) resource.TestCheckFunc {
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

		found, err := resources.RemoteAccessVpnUser.Get(rs.Primary.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Remote Access VPN User not found")
		}

		return nil
	}
}

func testAccCheckRemoteAccessVPNUserCreateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cca.CcaClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "cloudca_vpn_user" {
			if rs.Primary.Attributes["environment_id"] == "" {
				return fmt.Errorf("Environment ID is missing")
			}

			resources, err := getResourcesForEnvironmentID(client, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			_, err = resources.RemoteAccessVpnUser.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Remote Access VPN User still exists")
			}
		}
	}

	return nil
}
