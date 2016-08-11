package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceCloudcaTier() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaTierCreate,
		Read:   resourceCloudcaTierRead,
		Update: resourceCloudcaTierUpdate,
		Delete: resourceCloudcaTierDelete,

		Schema: map[string]*schema.Schema{
			"service_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A cloudca service code",
			},
			"environment_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of environment where tier should be created",
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
			"vpc": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or id of the VPC",
			},
			"network_offering": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The network offering name or id (e.g. "Standard Tier" or "Load Balanced Tier")`,
			},
			"network_acl": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ACL name or id",
			},
		},
	}
}

func resourceCloudcaTierCreate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	networkOfferingId, nerr := retrieveNetworkOfferingId(&ccaResources, d.Get("network_offering").(string))
	if nerr != nil {
		return nerr
	}

	vpcId, verr := retrieveVpcId(&ccaResources, d.Get("vpc").(string))
	if verr != nil {
		return verr
	}

	networkAclId, aerr := retrieveNetworkAclId(&ccaResources, d.Get("network_acl").(string))
	if aerr != nil {
		return aerr
	}

	tierToCreate := cloudca.Tier{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		VpcId:             vpcId,
		NetworkOfferingId: networkOfferingId,
		NetworkAclId:      networkAclId,
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
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

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

	vpc, vErr := ccaResources.Vpcs.Get(tier.VpcId)
	if vErr != nil {
		if ccaError, ok := vErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Vpc %s not found", tier.VpcId)
				d.SetId("")
				return nil
			}
		}
		return vErr
	}

	acl, aErr := ccaResources.NetworkAcls.Get(tier.NetworkAclId)
	if aErr != nil {
		if ccaError, ok := aErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("ACL %s not found", tier.NetworkAclId)
				d.SetId("")
				return nil
			}
		}
		return aErr
	}

	// Update the config
	d.Set("name", tier.Name)
	d.Set("description", tier.Description)
	setValueOrID(d, "network_offering", offering.Name, tier.NetworkOfferingId)
	setValueOrID(d, "vpc", vpc.Name, tier.VpcId)
	setValueOrID(d, "network_acl", acl.Name, tier.NetworkAclId)

	return nil
}

func resourceCloudcaTierUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	d.Partial(true)

	if d.HasChange("name") || d.HasChange("description") {
		newName := d.Get("name").(string)
		newDescription := d.Get("description").(string)
		_, err := ccaResources.Tiers.Update(d.Id(), cloudca.Tier{Id: d.Id(), Name: newName, Description: newDescription})
		if err != nil {
			return err
		}
	}

	if d.HasChange("network_acl") {
		networkAclId, aerr := retrieveNetworkAclId(&ccaResources, d.Get("network_acl").(string))
		if aerr != nil {
			return aerr
		}
		_, aclErr := ccaResources.Tiers.ChangeAcl(d.Id(), networkAclId)
		if aclErr != nil {
			return aclErr
		}
	}

	d.Partial(false)

	return nil
}

func resourceCloudcaTierDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)
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

func retrieveNetworkAclId(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}
	acls, err := ccaRes.NetworkAcls.List()
	if err != nil {
		return "", err
	}
	for _, acl := range acls {
		if strings.EqualFold(acl.Name, name) {
			log.Printf("Found network acl: %+v", acl)
			return acl.Id, nil
		}
	}
	return "", fmt.Errorf("Network ACL with name %s not found", name)
}
