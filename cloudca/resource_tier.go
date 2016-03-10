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
      Update: resourceCloudcaTierRead,
      Delete: resourceCloudcaTierRead,

      Schema: map[string]*schema.Schema{
         "service_code": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            ForceNew: true,
            Description: "A cloudca service code",
         },
         "environment_name": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            ForceNew: true,
            Description: "Name of environment where tier should be created",
         },
         "name": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            Description: "Name of tier",
         },
         "description": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            Description: "Description of tier",
         },
         "vpc": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            ForceNew: true,
            Description: "Name or id of the VPC",
            StateFunc: func(val interface{}) string {
               return strings.ToLower(val.(string))
            },
         },
         "network_offering": &schema.Schema{
            Type:     schema.TypeString,
            Required:true,
            ForceNew: true,
            Description: `The network offering name or id (e.g. "Standard Tier" or "Load Balanced Tier")`,
         },
         "network_acl": &schema.Schema{
            Type:     schema.TypeString,
            Required:true,
            Description: "The network ACL name or id",
         },
      },
   }
}

func resourceCloudcaTierCreate(d *schema.ResourceData, meta interface{}) error {
   ccaClient := meta.(*gocca.CcaClient)
   resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
   ccaResources := resources.(cloudca.Resources)

   networkOfferingId, nerr := retrieveNetworkOfferingId(&ccaResources, d.Get("network_offering").(string))
   if nerr != nil {
      return nerr
   }

   vpcId, verr := retrieveNetworkOfferingId(&ccaResources, d.Get("vpc").(string))
   if verr != nil {
      return verr
   }
   
   networkAclId, aerr := retrieveNetworkOfferingId(&ccaResources, d.Get("network_acl").(string))
   if aerr != nil {
      return aerr
   }

   tierToCreate := cloudca.Tier{
      Name: d.Get("name").(string),
      Description: d.Get("description").(string),
      VpcId: vpcId,
      NetworkOfferingId: networkOfferingId,
      NetworkAclId: networkAclId,
   }

   newTier, err := ccaResources.Tiers.Create(tierToCreate)
   if err != nil {
      return fmt.Errorf("Error creating the new tier %s: %s", tierToCreate.Name, err)
   }
   d.SetId(newTier.Id)

   return resourceCloudcaTierRead(d, meta)
}

func resourceCloudcaTierRead(d *schema.ResourceData, meta interface{}) error {
   ccaClient := meta.(*gocca.CcaClient)
   resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
   ccaResources := resources.(cloudca.Resources)

   // Get the vpc details
   vpc, err := ccaResources.Tiers.Get(d.Id())
   if err != nil {
      if ccaError, ok := err.(api.CcaErrorResponse); ok {
         if ccaError.StatusCode == 404 {
            fmt.Errorf("VPC %s does no longer exist", d.Get("name").(string))
            d.SetId("")
            return nil
         }
      }
      return err
   }

   vpcOffering, offErr := ccaResources.VpcOfferings.Get(vpc.NetworkOfferingId)
   if offErr != nil {
      if ccaError, ok := offErr.(api.CcaErrorResponse); ok {
         if ccaError.StatusCode == 404 {
            fmt.Errorf("VPC offering id=%s does no longer exist", vpc.NetworkOfferingId)
            d.SetId("")
            return nil
         }
      }
      return offErr
   }

   // Update the config
   d.Set("name", vpc.Name)
   d.Set("description", vpc.Description)
   setValueOrID(d, "vpc_offering", strings.ToLower(vpcOffering.Name), vpc.NetworkOfferingId)
   d.Set("network_domain", vpc.NetworkAclId)

   return nil
}

// func resourceCloudcaTierUpdate(d *schema.ResourceData, meta interface{}) error {
//    ccaClient := meta.(*gocca.CcaClient)
//    resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
//    ccaResources := resources.(cloudca.Resources)

//    if d.HasChange("name") || d.HasChange("description") {
//       newName := d.Get("name").(string)
//       newDescription := d.Get("name").(string)
//       log.Printf("[DEBUG] Details have changed updating VPC.....")
//       _, err := ccaResources.Tiers.Update(d.Id(), cloudca.Tier{Id: d.Id(), Name: newName, Description: newDescription})
//       if err != nil {
//          return err
//       }
//    }

//    return nil
// }

// func resourceCloudcaTierDelete(d *schema.ResourceData, meta interface{}) error {
//    ccaClient := meta.(*gocca.CcaClient)
//    resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
//    ccaResources := resources.(cloudca.Resources)

//    fmt.Println("[INFO] Destroying VPC: %s", d.Get("name").(string))
//    if _, err := ccaResources.Tiers.Destroy(d.Id()); err != nil {
//       if ccaError, ok := err.(api.CcaErrorResponse); ok {
//          if ccaError.StatusCode == 404 {
//             fmt.Errorf("VPC %s does no longer exist", d.Get("name").(string))
//             d.SetId("")
//             return nil
//          }
//       }
//       return err
//    }

//    return nil
// }

func retrieveNetworkOfferingId(ccaRes *cloudca.Resources, name string) (id string, err error) {
   if isID(name) {
      return name, nil
   }
   offerings, err := ccaRes.NetworkOfferings.List()
   if err != nil {
      return "", err
   }
   for _, offering := range offerings {
       if strings.EqualFold(offering.Name,name) {
         log.Printf("Found network offering: %+v", offering)
         return offering.Id, nil
       }
   }
   return "", fmt.Errorf("Network offering with name %s not found", name)
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
       if strings.EqualFold(vpc.Name,name) {
         log.Printf("Found vpc: %+v", vpc)
         return vpc.Id, nil
       }
   }
   return "", fmt.Errorf("Vpc with name %s not found", name)
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
       if strings.EqualFold(acl.Name,name) {
         log.Printf("Found network acl: %+v", acl)
         return acl.Id, nil
       }
   }
   return "", fmt.Errorf("Network ACL with name %s not found", name)
}
