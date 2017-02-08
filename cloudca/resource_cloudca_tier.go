package cloudca

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaTier() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaTierCreate,
		Read:   resourceCloudcaTierRead,
		Update: resourceCloudcaTierUpdate,
		Delete: resourceCloudcaTierDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where tier should be created",
			},
			"organization_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Entry point of organization",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of tier",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of tier",
			},
			"vpc_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC",
			},
			"network_offering": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The network offering name or id (e.g. "Standard Tier" or "Load Balanced Tier")`,
			},
			"network_acl_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the network ACL",
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudcaTierCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	networkOfferingId, nerr := retrieveNetworkOfferingId(&ccaResources, d.Get("network_offering").(string))
	if nerr != nil {
		return nerr
	}

	tierToCreate := cloudca.Tier{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		VpcId:             d.Get("vpc_id").(string),
		NetworkOfferingId: networkOfferingId,
		NetworkAclId:      d.Get("network_acl_id").(string),
	}
	options := map[string]string{}
	if orgId, ok := d.GetOk("organization_code"); ok {
		options["org_id"] = orgId.(string)
	}
	newTier, err := ccaResources.Tiers.Create(tierToCreate, options)
	if err != nil {
		return fmt.Errorf("Error creating the new tier %s: %s", tierToCreate.Name, err)
	}
	d.SetId(newTier.Id)
	return resourceCloudcaTierRead(d, meta)
}

func resourceCloudcaTierRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	tier, err := ccaResources.Tiers.Get(d.Id())
	if err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				log.Printf("Tier %s was not found", d.Get("name").(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}

	offering, offErr := ccaResources.NetworkOfferings.Get(tier.NetworkOfferingId)
	if offErr != nil {
		if ccaError, ok := offErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Network offering %s not found", tier.NetworkOfferingId)
				d.SetId("")
				return nil
			}
		}
		return offErr
	}

	// Update the config
	d.Set("name", tier.Name)
	d.Set("description", tier.Description)
	setValueOrID(d, "network_offering", offering.Name, tier.NetworkOfferingId)
	d.Set("vpc_id", tier.VpcId)
	d.Set("network_acl_id", tier.NetworkAclId)
	d.Set("cidr", tier.Cidr)
	return nil
}

func resourceCloudcaTierUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	d.Partial(true)

	if d.HasChange("name") || d.HasChange("description") {
		newName := d.Get("name").(string)
		newDescription := d.Get("description").(string)
		_, err := ccaResources.Tiers.Update(d.Id(), cloudca.Tier{Id: d.Id(), Name: newName, Description: newDescription})
		if err != nil {
			return err
		}
	}

	if d.HasChange("network_acl_id") {
		_, aclErr := ccaResources.Tiers.ChangeAcl(d.Id(), d.Get("network_acl_id").(string))
		if aclErr != nil {
			return aclErr
		}
	}

	d.Partial(false)

	return nil
}

func resourceCloudcaTierDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	if _, err := ccaResources.Tiers.Delete(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Tier %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	return nil
}

func retrieveNetworkOfferingId(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}
	offerings, err := ccaRes.NetworkOfferings.List()
	if err != nil {
		return "", err
	}
	for _, offering := range offerings {
		if strings.EqualFold(offering.Name, name) {
			log.Printf("Found network offering: %+v", offering)
			return offering.Id, nil
		}
	}
	return "", fmt.Errorf("Network offering with name %s not found", name)
}
