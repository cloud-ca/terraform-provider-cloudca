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

func resourceCloudcaPublicIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaPublicIpCreate,
		Read:   resourceCloudcaPublicIpRead,
		Update: resourceCloudcaPublicIpUpdate,
		Delete: resourceCloudcaPublicIpDelete,

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
				Description: "Name of environment where the public IP should be created",
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or id of the VPC",
			},
		},
	}
}

func resourceCloudcaPublicIpCreate(d *schema.ResourceData, meta interface{}) error {
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

	publicIpToCreate := cloudca.PublicIp{
		VpcId:             vpcId,
	}
	newPublicIp, err := ccaResources.PublicIps.Acquire(publicIpToCreate)
	if err != nil {
		return fmt.Errorf("Error acquiring the new public ip %s", err)
	}
	d.SetId(newPublicIp.Id)
	return resourceCloudcaPublicIpRead(d, meta)
}

func resourceCloudcaPublicIpRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	publicIp, err := ccaResources.PublicIp.Get(d.Id())
	if err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				log.Printf("Public Ip with id='%s' was not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	vpc, vErr := ccaResources.Vpcs.Get(publicIp.VpcId)
	if vErr != nil {
		if ccaError, ok := vErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Vpc %s not found", publicIp.VpcId)
				d.SetId("")
				return nil
			}
		}
		return vErr
	}

	setValueOrID(d, "vpc", vpc.Name, publicIp.VpcId)

	return nil
}

func resourceCloudcaPublicIpUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudcaPublicIpDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)
	if _, err := ccaResources.PublicIp.Release(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Public Ip %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	return nil
}

func retrieveVpcId(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}
	vpcs, err := ccaRes.Vpcs.List()
	if err != nil {
		return "", err
	}
	for _, vpc := range vpcs {
		if strings.EqualFold(vpc.Name, name) {
			log.Printf("Found vpc: %+v", vpc)
			return vpc.Id, nil
		}
	}
	return "", fmt.Errorf("Vpc with name %s not found", name)
}