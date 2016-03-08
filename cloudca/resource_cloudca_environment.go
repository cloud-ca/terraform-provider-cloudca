package cloudca

import (
   "fmt"
   "log"
   "strings"
   "github.com/cloud-ca/go-cloudca" 
   "github.com/cloud-ca/go-cloudca/api" 
   "github.com/cloud-ca/go-cloudca/configuration"
   "github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaEnvironment() *schema.Resource {
   return &schema.Resource{
      Create: resourceCloudcaEnvironmentCreate,
      Read:   resourceCloudcaEnvironmentRead,
      Update: resourceCloudcaEnvironmentRead,
      Delete: resourceCloudcaEnvironmentDelete,

      Schema: map[string]*schema.Schema{
         "organization_code": &schema.Schema{
            Type:     schema.TypeString,
            ForceNew: true,
            Required: true,
            Description: "Organization's entry point, i.e. <entry_point>.cloud.ca",
            StateFunc: func(val interface{}) string {
               return strings.ToLower(val.(string))
            },
         },
         "service_code": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            ForceNew: true,
            Description: "A cloudca service code",
         },
         "name": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            Description: "Name of environment to be created. Must be lower case, contain alphanumeric charaters, underscores or dashes",
         },
         "description": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
            Description: "Description for the environment",
         },
         "membership": &schema.Schema{
            Type:     schema.TypeString,
            Optional:true,
            Description: "Environment membership for users (options: all, many)",
         },
         "admin_role": &schema.Schema{
            Type:     schema.TypeList,
            Elem:     &schema.Schema{Type: schema.TypeString},
            Optional:true,
            Description: "List of users that will be given Environment Admin role",
         },
         "user_role": &schema.Schema{
            Type:     schema.TypeList,
            Elem:     &schema.Schema{Type: schema.TypeString},
            Optional:true,
            Description: "List of users that will be given User role",
         },
         "read_only_role": &schema.Schema{
            Type:     schema.TypeList,
            Elem:     &schema.Schema{Type: schema.TypeString},
            Optional:true,
            Description: "List of users that will be given Read-only role",
         },
      },
   }
}

func resourceCloudcaEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
   ccaClient := meta.(*gocca.CcaClient)

   organizationId, oerr := getOrganizationId(ccaClient, d.Get("organization_code").(string))
   if oerr != nil {
      return oerr
   }

   connectionId, cerr := getServiceConnectionId(ccaClient, d.Get("service_code").(string))

   if cerr != nil {
      return cerr
   }

   environmentToCreate := configuration.Environment{
      Name: d.Get("name").(string),
      Description: d.Get("description").(string),
      ServiceConnection: configuration.ServiceConnection{Id: connectionId,},
   }

   if organizationId != "" {
      environmentToCreate.Organization = configuration.Organization{Id:organizationId,}
   }

   if membership, ok := d.GetOk("membership"); ok {
      environmentToCreate.Membership = membership.(string)
   }

   _, adminRoleExists := d.GetOk("admin_role")
   _, userRoleExists := d.GetOk("user_role")
   _, readOnlyRoleExists := d.GetOk("read_only_role")

   if adminRoleExists || userRoleExists || readOnlyRoleExists {
      
      users, uerr := ccaClient.Users.ListWithOptions(map[string]string{"tenantId":organizationId})
      if uerr != nil{
         return uerr
      }

      environmentToCreate.Roles = []configuration.Role{}

      if adminRoleExists {
         role, err := mapUsersToRole("Environment admin", d.Get("admin_role").([]interface{}), users)
         if err != nil{
            return err
         }
         environmentToCreate.Roles = append(environmentToCreate.Roles, role)
      }

      if userRoleExists {
         role, err := mapUsersToRole("User", d.Get("user_role").([]interface{}), users)
         if err != nil{
            return err
         }
         environmentToCreate.Roles = append(environmentToCreate.Roles, role)
      }

      if readOnlyRoleExists {
         role, err := mapUsersToRole("Read-only", d.Get("read_only_role").([]interface{}), users)
         if err != nil{
            return err
         }
         environmentToCreate.Roles = append(environmentToCreate.Roles, role)
      }
   }

   newEnvironment, err := ccaClient.Environments.Create(environmentToCreate)
   if err != nil {
      return fmt.Errorf("Error creating the new environment %s: %s", environmentToCreate.Name, err)
   }

   d.SetId(newEnvironment.Id)

   return resourceCloudcaEnvironmentRead(d, meta)
}

func mapUsersToRole(roleName string, roleUserList []interface{}, users []configuration.User) (configuration.Role, error) {
   role := configuration.Role{
      Name: roleName,
      Users:[]configuration.User{},
   }

   for _, userToFind := range roleUserList {
      found := false
      for _, user := range users{
         if strings.EqualFold(user.Username,userToFind.(string)) {
            found = true
            role.Users = append(role.Users, configuration.User{Id:user.Id,})
            break
         }
      }
      if !found {
         return configuration.Role{},fmt.Errorf("User %s was not found", userToFind)
      }
   }
   return role, nil
}

func resourceCloudcaEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
   ccaClient := meta.(*gocca.CcaClient)
   _, err := ccaClient.Environments.Get(d.Id())
   if err != nil {
      if ccaError, ok := err.(api.CcaErrorResponse); ok {
         if ccaError.StatusCode == 404 {
            fmt.Errorf("Environment %s does not exist", d.Id())
            d.SetId("")
            return nil
         }
      }
      return err
   }
   return nil
}

// func resourceCloudcaEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
//    ccaClient := meta.(*gocca.CcaClient)
//    resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
//    ccaResources := resources.(cloudca.Resources)

//    d.Partial(true)

//    if d.HasChange("compute_offering") {
//       newComputeOffering := d.Get("compute_offering").(string)
//       log.Printf("[DEBUG] Compute offering has changed for %s, changing compute offering...", newComputeOffering)
//       newComputeOfferingId, ferr := retrieveComputeOfferingID(&ccaResources, newComputeOffering)
//       if ferr != nil {
//          return ferr
//       }
//       _, err := ccaResources.Instances.ChangeComputeOffering(d.Id(), newComputeOfferingId)
//       if err != nil {
//          return err
//       }
//       d.SetPartial("compute_offering")
//    }

//    if d.HasChange("ssh_key_name") {
//       sshKeyName := d.Get("ssh_key_name").(string)
//       log.Printf("[DEBUG] SSH key name has changed for %s, associating new SSH key...", sshKeyName)
//       _, err := ccaResources.Instances.AssociateSSHKey(d.Id(), sshKeyName)
//       if err != nil {
//          return err
//       }
//       d.SetPartial("ssh_key_name")
//    }

//    d.Partial(false)

//    return nil
// }

func resourceCloudcaEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
   ccaClient := meta.(*gocca.CcaClient)
   fmt.Println("[INFO] Destroying environment: %s", d.Get("name").(string))
   if _, err := ccaClient.Environments.Delete(d.Id()); err != nil {
      if ccaError, ok := err.(api.CcaErrorResponse); ok {
         if ccaError.StatusCode == 404 {
            fmt.Errorf("Environment %s does not exist", d.Get("name").(string))
            d.SetId("")
            return nil
         }
      }
      return err
   }
   return nil
}

func getServiceConnectionId(ccaClient *gocca.CcaClient, serviceCode string) (id string, err error){
   if isID(serviceCode){
      return serviceCode, nil
   }
   connections, cerr := ccaClient.ServiceConnections.List();
   if cerr != nil {
      return "", cerr
   }
   for _, connection := range connections {
      if strings.EqualFold(connection.ServiceCode,serviceCode) {
         log.Printf("Found service connection : %+v", connection)
         return connection.Id, nil
      }
   }
   return "", nil
}


func getOrganizationId(ccaClient *gocca.CcaClient, entryPoint string) (id string, err error) {
   if isID(entryPoint){
      return entryPoint, nil
   }
   orgs, err := ccaClient.Organizations.List()
   if err != nil {
      return "", err
   }
   for _, org := range orgs {
       if strings.EqualFold(org.EntryPoint,entryPoint) {
         log.Printf("Found organization: %+v", org)
         return org.Id, nil
       }
   }
   return "", fmt.Errorf("Organization with entry point %s not found", entryPoint)
}