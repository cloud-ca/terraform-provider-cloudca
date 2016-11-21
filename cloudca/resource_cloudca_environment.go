package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/configuration"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

const (
	//role name
	ENVIRONMENT_ADMIN_ROLE = "Environment admin"
	USER_ROLE              = "User"
	READ_ONLY_ROLE         = "Read-only"

	//fields
	ORGANIZATION_CODE    = "organization_code"
	SERVICE_CODE         = "service_code"
	NAME                 = "name"
	DESCRIPTION          = "description"
	ADMIN_ROLE_USERS     = "admin_role"
	USER_ROLE_USERS      = "user_role"
	READ_ONLY_ROLE_USERS = "read_only_role"
)

func resourceCloudcaEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaEnvironmentCreate,
		Read:   resourceCloudcaEnvironmentRead,
		Update: resourceCloudcaEnvironmentUpdate,
		Delete: resourceCloudcaEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			ORGANIZATION_CODE: &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Organization's entry point, i.e. <entry_point>.cloud.ca",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			SERVICE_CODE: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A cloudca service code",
			},
			NAME: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of environment to be created. Must be lower case, contain alphanumeric charaters, underscores or dashes",
			},
			DESCRIPTION: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description for the environment",
			},
			ADMIN_ROLE_USERS: &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of users that will be given Environment Admin role",
			},
			USER_ROLE_USERS: &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of users that will be given User role",
			},
			READ_ONLY_ROLE_USERS: &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of users that will be given Read-only role",
			},
		},
	}
}

func resourceCloudcaEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	environment, err := ccaClient.Environments.Get(d.Id())
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

	adminRoleUsers, userRoleUsers, readOnlyRoleUsers := getUsersFromRoles(environment)
	adminRole, _ := d.GetOk(ADMIN_ROLE_USERS)
	userRole, _ := d.GetOk(USER_ROLE_USERS)
	readOnlyRole, _ := d.GetOk(READ_ONLY_ROLE_USERS)

	d.Set(NAME, environment.Name)
	d.Set(DESCRIPTION, environment.Description)
	d.Set(ADMIN_ROLE_USERS, getListOfUsersByIdOrUsername(adminRoleUsers, adminRole.(*schema.Set)))
	d.Set(USER_ROLE_USERS, getListOfUsersByIdOrUsername(userRoleUsers, userRole.(*schema.Set)))
	d.Set(READ_ONLY_ROLE_USERS, getListOfUsersByIdOrUsername(readOnlyRoleUsers, readOnlyRole.(*schema.Set)))

	return nil
}

func resourceCloudcaEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)

	environment, err := getEnvironmentFromConfig(ccaClient, d)
	if err != nil {
		return fmt.Errorf("Error parsing environment %s: %s", environment.Name, err)
	}

	newEnvironment, newErr := ccaClient.Environments.Create(*environment)
	if newErr != nil {
		return fmt.Errorf("Error creating the new environment %s: %s", environment.Name, newErr)
	}

	d.SetId(newEnvironment.Id)

	return resourceCloudcaEnvironmentRead(d, meta)
}

func resourceCloudcaEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	environment, err := getEnvironmentFromConfig(ccaClient, d)
	if err != nil {
		return fmt.Errorf("Error parsing environment %s: %s", environment.Name, err)
	}
	_, uerr := ccaClient.Environments.Update(d.Id(), *environment)
	if uerr != nil {
		return fmt.Errorf("Error updating environment %s: %s", environment.Name, uerr)
	}
	return resourceCloudcaEnvironmentRead(d, meta)
}

func resourceCloudcaEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	fmt.Println("[INFO] Destroying environment: %s", d.Get(NAME).(string))
	if _, err := ccaClient.Environments.Delete(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Environment %s does not exist", d.Get(NAME).(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}
	return nil
}

func getEnvironmentFromConfig(ccaClient *cca.CcaClient, d *schema.ResourceData) (*configuration.Environment, error) {
	environment := configuration.Environment{}
	environment.Name = d.Get(NAME).(string)
	environment.Description = d.Get(DESCRIPTION).(string)

	organizationId, oerr := getOrganizationId(ccaClient, d.Get(ORGANIZATION_CODE).(string))
	if oerr != nil {
		return &environment, oerr
	}

	connectionId, cerr := getServiceConnectionId(ccaClient, d.Get(SERVICE_CODE).(string))
	if cerr != nil {
		return &environment, cerr
	}

	environment.Organization = configuration.Organization{Id: organizationId}
	environment.ServiceConnection = configuration.ServiceConnection{Id: connectionId}

	adminRole, adminRoleExists := d.GetOk(ADMIN_ROLE_USERS)
	userRole, userRoleExists := d.GetOk(USER_ROLE_USERS)
	readOnlyRole, readOnlyRoleExists := d.GetOk(READ_ONLY_ROLE_USERS)

	if adminRoleExists || userRoleExists || readOnlyRoleExists {

		users, uerr := ccaClient.Users.ListWithOptions(map[string]string{"organizationId": organizationId})
		if uerr != nil {
			return &environment, uerr
		}

		environment.Roles = []configuration.Role{}

		if adminRoleExists {
			role, err := mapUsersToRole(ENVIRONMENT_ADMIN_ROLE, adminRole.(*schema.Set).List(), users)
			if err != nil {
				return &environment, err
			}
			environment.Roles = append(environment.Roles, role)
		}

		if userRoleExists {
			role, err := mapUsersToRole(USER_ROLE, userRole.(*schema.Set).List(), users)
			if err != nil {
				return &environment, err
			}
			environment.Roles = append(environment.Roles, role)
		}

		if readOnlyRoleExists {
			role, err := mapUsersToRole(READ_ONLY_ROLE, readOnlyRole.(*schema.Set).List(), users)
			if err != nil {
				return &environment, err
			}
			environment.Roles = append(environment.Roles, role)
		}
	}
	return &environment, nil
}

func getListOfUsersByIdOrUsername(roleUsers []configuration.User, usersWithIdOrName *schema.Set) *schema.Set {
	mappedList := []interface{}{}
	for _, user := range roleUsers {
		found := false
		for _, idOrUsername := range usersWithIdOrName.List() {
			if isID(idOrUsername.(string)) {
				if strings.EqualFold(user.Id, idOrUsername.(string)) {
					found = true
					mappedList = append(mappedList, user.Id)
					break
				}
			} else if strings.EqualFold(user.Username, idOrUsername.(string)) {
				found = true
				mappedList = append(mappedList, user.Username)
				break
			}
		}
		if !found {
			mappedList = append(mappedList, user.Username)
		}
	}
	return schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeString}), mappedList)
}

func getUsersFromRoles(environment *configuration.Environment) (adminRoleUsers []configuration.User, userRoleUsers []configuration.User, readOnlyRoleUsers []configuration.User) {
	for _, envRole := range environment.Roles {
		switch {
		case strings.EqualFold(envRole.Name, ENVIRONMENT_ADMIN_ROLE):
			for _, user := range envRole.Users {
				adminRoleUsers = append(adminRoleUsers, user)
			}
		case strings.EqualFold(envRole.Name, USER_ROLE):
			for _, user := range envRole.Users {
				userRoleUsers = append(userRoleUsers, user)
			}
		case strings.EqualFold(envRole.Name, READ_ONLY_ROLE):
			for _, user := range envRole.Users {
				readOnlyRoleUsers = append(readOnlyRoleUsers, user)
			}
		}
	}
	return
}

func mapUsersToRole(roleName string, userList []interface{}, users []configuration.User) (configuration.Role, error) {
	role := configuration.Role{
		Name:  roleName,
		Users: []configuration.User{},
	}

	for _, userToFind := range userList {
		if isID(userToFind.(string)) {
			role.Users = append(role.Users, configuration.User{Id: userToFind.(string)})
			continue
		}
		found := false
		for _, user := range users {
			if strings.EqualFold(user.Username, userToFind.(string)) {
				found = true
				role.Users = append(role.Users, configuration.User{Id: user.Id})
				break
			}
		}
		if !found {
			return configuration.Role{}, fmt.Errorf("User %s was not found", userToFind)
		}
	}
	return role, nil
}

func getServiceConnectionId(ccaClient *cca.CcaClient, serviceCode string) (id string, err error) {
	if isID(serviceCode) {
		return serviceCode, nil
	}
	connections, cerr := ccaClient.ServiceConnections.List()
	if cerr != nil {
		return "", cerr
	}
	for _, connection := range connections {
		if strings.EqualFold(connection.ServiceCode, serviceCode) {
			log.Printf("Found service connection : %+v", connection)
			return connection.Id, nil
		}
	}
	return "", nil
}

func getOrganizationId(ccaClient *cca.CcaClient, entryPoint string) (id string, err error) {
	if isID(entryPoint) {
		return entryPoint, nil
	}
	orgs, err := ccaClient.Organizations.List()
	if err != nil {
		return "", err
	}
	for _, org := range orgs {
		if strings.EqualFold(org.EntryPoint, entryPoint) {
			log.Printf("Found organization: %+v", org)
			return org.Id, nil
		}
	}
	return "", fmt.Errorf("Organization with entry point %s not found", entryPoint)
}
