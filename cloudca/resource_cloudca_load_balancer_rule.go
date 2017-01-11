package cloudca

import (
   "fmt"
   "github.com/cloud-ca/go-cloudca"
   "github.com/cloud-ca/go-cloudca/api"
   "github.com/cloud-ca/go-cloudca/services/cloudca"
   "github.com/hashicorp/terraform/helper/schema"
   "strconv"
   "errors"
)

func resourceCloudcaLoadBalancerRule() *schema.Resource {
   return &schema.Resource{
      Create: createLbr,
      Read:   readLbr,
      Delete: deleteLbr,
      Update: updateLbr,

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
            Description: "Name of environment where load balancer rule should be created",
         },
         "name": &schema.Schema{
            Type:        schema.TypeString,
            Required:    true,
            Description: "Name of the load balancer rule",
         },
         "public_ip_id": &schema.Schema{
            Type:        schema.TypeString,
            Required:    true,
            ForceNew:    true,
            Description: "The public IP to which the rule should be applied",
         },
         "public_ip": &schema.Schema{
            Type:        schema.TypeString,
            Computed:    true,
         },
         "network_id": &schema.Schema{
            Type:        schema.TypeString,
            Optional:    true,
            ForceNew:    true,
            Computed: true,
            Description: "The network ID to bind to",
         },
         "protocol": &schema.Schema{
            Type:        schema.TypeString,
            Required:    true,
            ForceNew:    true,
            Description: "The protocol that this rule should use (eg. TCP, UDP)",
         },
         "algorithm": &schema.Schema{
            Type:        schema.TypeString,
            Required:    true,
            Description: "The algorithm used to load balance",
         },
         "public_port": &schema.Schema{
            Type:        schema.TypeInt,
            Required:    true,
            ForceNew:    true,
            Description: "The port on the public IP",
         },
         "private_port": &schema.Schema{
            Type:        schema.TypeInt,
            Required:    true,
            ForceNew:    true,
            Description: "The port to which the traffic will be load balanced internally",
         },
         "instance_ids": &schema.Schema{
            Type:        schema.TypeList,
            Optional:    true,
            Description: "List of instance ids that will be load balanced",
            Elem:        &schema.Schema{Type: schema.TypeString},
         },
         "stickiness_method": &schema.Schema{
            Type:        schema.TypeString,
            Optional:    true,
            Description: "The stickiness method",
         },
         "stickiness_params": &schema.Schema{
            Type:        schema.TypeMap,
            Optional:    true,
            Description: "The stickiness policy parameters",
         },
      },
   }
}

func createLbr(d *schema.ResourceData, meta interface{}) error {
   client := meta.(*cca.CcaClient)
   resources, err := client.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
   ccaResources := resources.(cloudca.Resources)

   if err != nil {
      return err
   }

   lbr := cloudca.LoadBalancerRule{
      Name:             d.Get("name").(string),
      PublicIpId:       d.Get("public_ip_id").(string),
      NetworkId:        d.Get("network_id").(string),
      Protocol:         d.Get("protocol").(string),
      Algorithm:        d.Get("algorithm").(string),
      PublicPort:       strconv.Itoa(d.Get("public_port").(int)),
      PrivatePort:      strconv.Itoa(d.Get("private_port").(int)),
   }

   _, instanceIdsPresent := d.GetOk("instance_ids")

   if instanceIdsPresent {
      var instanceIds []string
      for _, id := range d.Get("instance_ids").([]interface{}) {
         instanceIds = append(instanceIds, id.(string))
      }
      lbr.InstanceIds = instanceIds
   }

   
   if stickinessMethod, ok := d.GetOk("stickiness_method"); ok {
      lbr.StickinessMethod = stickinessMethod.(string)
   }

   
   if stickinessParams, ok := d.GetOk("stickiness_params"); ok {
      lbr.StickinessPolicyParameters = getStickinessPolicyParameterMap(stickinessParams.(map[string]interface{}))
   }

   newLbr, err := ccaResources.LoadBalancerRules.Create(lbr)
   if err != nil {
      return err
   }

   d.SetId(newLbr.Id)
   return readLbr(d, meta)
}

func readLbr(d *schema.ResourceData, meta interface{}) error {
   client := meta.(*cca.CcaClient)
   resources, _ := client.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
   ccaResources := resources.(cloudca.Resources)

   lbr, err := ccaResources.LoadBalancerRules.Get(d.Id())
   if err != nil {
      return handleLbrNotFoundError(err, d)
   }

   d.Set("name", lbr.Name)
   d.Set("public_ip_id", lbr.PublicIpId)
   d.Set("network_id", lbr.NetworkId)
   d.Set("instance_ids", lbr.InstanceIds)
   d.Set("algorithm", lbr.Algorithm)
   d.Set("protocol", lbr.Protocol)
   d.Set("public_port", lbr.PublicPort)
   d.Set("private_port", lbr.PrivatePort)
   d.Set("public_ip", lbr.PublicIp)
   d.Set("stickiness_method", lbr.StickinessMethod)
   d.Set("stickiness_params", lbr.StickinessPolicyParameters)

   return nil
}

func deleteLbr(d *schema.ResourceData, meta interface{}) error {
   client := meta.(*cca.CcaClient)
   resources, _ := client.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
   ccaResources := resources.(cloudca.Resources)

   if err := ccaResources.LoadBalancerRules.Delete(d.Id()); err != nil {
      return handleLbrNotFoundError(err, d)
   }
   return nil
}

func updateLbr(d *schema.ResourceData, meta interface{}) error {
   client := meta.(*cca.CcaClient)
   resources, err := client.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
   ccaResources := resources.(cloudca.Resources)

   if err != nil {
      return err
   }

   d.Partial(true)

   if d.HasChange("stickiness_method") || d.HasChange("stickiness_params") {
      stickinessMethod := d.Get("stickiness_method")
      if len(stickinessMethod.(string)) > 0 {
         var stickinessPolicyParameters map[string]string
         if stickinessParams, ok := d.GetOk("stickiness_params"); ok {
            stickinessPolicyParameters = getStickinessPolicyParameterMap(stickinessParams.(map[string]interface{}))
         }
         err := ccaResources.LoadBalancerRules.SetLoadBalancerRuleStickinessPolicy(d.Id(), stickinessMethod.(string), stickinessPolicyParameters)
         if err != nil {
            return err
         }
      }else{
         
         if _, ok := d.GetOk("stickiness_params"); ok { 
            return errors.New("Stickiness params should be removed if the stickiness method is removed")
         }
         err := ccaResources.LoadBalancerRules.RemoveLoadBalancerRuleStickinessPolicy(d.Id())
         if err != nil {
            return err
         }
      }
   }

   if d.HasChange("name") || d.HasChange("algorithm") {
      newName := d.Get("name").(string)
      newAlgorithm := d.Get("algorithm").(string)
      _, err := ccaResources.LoadBalancerRules.Update(cloudca.LoadBalancerRule{Id: d.Id(), Name: newName, Algorithm: newAlgorithm})
      if err != nil {
         return err
      }
   }

   if d.HasChange("instance_ids") {
      var instanceIds []string
      for _, id := range d.Get("instance_ids").([]interface{}) {
         instanceIds = append(instanceIds, id.(string))
      }

      instanceErr := ccaResources.LoadBalancerRules.SetLoadBalancerRuleInstances(d.Id(), instanceIds)
      if instanceErr != nil {
         return instanceErr
      }
   }
   d.Partial(false)
   return readLbr(d, meta)
}

func getStickinessPolicyParameterMap(policyMap map[string]interface{}) map[string]string {
   var paramMap = make(map[string]string)
   for k, v := range policyMap {
      paramMap[k] = v.(string)
   }
   return paramMap
}

func handleLbrNotFoundError(err error, d *schema.ResourceData) error {
   if ccaError, ok := err.(api.CcaErrorResponse); ok {
      if ccaError.StatusCode == 404 {
         fmt.Errorf("Load balancer rule with id %s was not found", d.Id())
         d.SetId("")
         return nil
      }
   }

   return err
}

