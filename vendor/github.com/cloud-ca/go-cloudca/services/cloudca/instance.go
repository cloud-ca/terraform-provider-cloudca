package cloudca

import (
	"encoding/json"
	"strings"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

const (
	INSTANCE_STATE_RUNNING = "Running"
	INSTANCE_STATE_STOPPED = "Stopped"
)

const (
	INSTANCE_START_OPERATION                   = "start"
	INSTANCE_STOP_OPERATION                    = "stop"
	INSTANCE_REBOOT_OPERATION                  = "reboot"
	INSTANCE_RECOVER_OPERATION                 = "recover"
	INSTANCE_PURGE_OPERATION                   = "purge"
	INSTANCE_RESET_PASSWORD_OPERATION          = "resetPassword"
	INSTANCE_CREATE_RECOVERY_POINT_OPERATION   = "createRecoveryPoint"
	INSTANCE_CHANGE_COMPUTE_OFFERING_OPERATION = "changeComputeOffering"
	INSTANCE_ASSOCIATE_SSH_KEY_OPERATION       = "associateSSHKey"
)

type Instance struct {
	Id                       string        `json:"id,omitempty"`
	Name                     string        `json:"name,omitempty"`
	State                    string        `json:"state,omitempty"`
	TemplateId               string        `json:"templateId,omitempty"`
	TemplateName             string        `json:"templateName,omitempty"`
	IsPasswordEnabled        bool          `json:"isPasswordEnabled,omitempty"`
	IsSSHKeyEnabled          bool          `json:"isSshKeyEnabled,omitempty"`
	Username                 string        `json:"username,omitempty"`
	Password                 string        `json:"password,omitempty"`
	SSHKeyName               string        `json:"sshKeyName,omitempty"`
	ComputeOfferingId        string        `json:"computeOfferingId,omitempty"`
	ComputeOfferingName      string        `json:"computeOfferingName,omitempty"`
	NewComputeOfferingId     string        `json:"newComputeOfferingId,omitempty"`
	CpuCount                 int           `json:"cpuCount,omitempty"`
	MemoryInMB               int           `json:"memoryInMB,omitempty"`
	ZoneId                   string        `json:"zoneId,omitempty"`
	ZoneName                 string        `json:"zoneName,omitempty"`
	ProjectId                string        `json:"projectId,omitempty"`
	NetworkId                string        `json:"networkId,omitempty"`
	NetworkName              string        `json:"networkName,omitempty"`
	VpcId                    string        `json:"vpcId,omitempty"`
	VpcName                  string        `json:"vpcName,omitempty"`
	MacAddress               string        `json:"macAddress,omitempty"`
	UserData                 string        `json:"userData,omitempty"`
	RecoveryPoint            RecoveryPoint `json:"recoveryPoint,omitempty"`
	IpAddress                string        `json:"ipAddress,omitempty"`
	IpAddressId              string        `json:"ipAddressId,omitempty"`
	PublicIps                []PublicIp    `json:"publicIPs,omitempty"`
	PublicKey                string        `json:"publicKey,omitempty"`
	AdditionalDiskOfferingId string        `json:"diskOfferingId,omitempty"`
	AdditionalDiskSizeInGb   string        `json:"additionalDiskSizeInGb,omitempty"`
	AdditionalDiskIops       string        `json:"additionalDiskIops,omitempty"`
	VolumeIdToAttach         string        `json:"volumeIdToAttach,omitempty"`
	PortsToForward           []string      `json:"portsToForward,omitempty"`
	RootVolumeSizeInGb       int           `json:"rootVolumeSizeInGb,omitempty"`
	DedicatedGroupId         string        `json:"dedicatedGroupId,omitempty"`
	AffinityGroupIds         []string      `json:"affinityGroupIds,omitempty"`
}

type DestroyOptions struct {
	PurgeImmediately     bool     `json:"purgeImmediately,omitempty"`
	DeleteSnapshots      bool     `json:"deleteSnapshots,omitempty"`
	PublicIpIdsToRelease []string `json:"publicIpIdsToRelease,omitempty"`
	VolumeIdsToDelete    []string `json:"volumeIdsToDelete,omitempty"`
}

func (instance *Instance) IsRunning() bool {
	return strings.EqualFold(instance.State, INSTANCE_STATE_RUNNING)
}

func (instance *Instance) IsStopped() bool {
	return strings.EqualFold(instance.State, INSTANCE_STATE_STOPPED)
}

type InstanceService interface {
	Get(id string) (*Instance, error)
	List() ([]Instance, error)
	ListWithOptions(options map[string]string) ([]Instance, error)
	Create(Instance) (*Instance, error)
	Destroy(id string, purge bool) (bool, error)
	DestroyWithOptions(id string, options DestroyOptions) (bool, error)
	Purge(id string) (bool, error)
	Recover(id string) (bool, error)
	Exists(id string) (bool, error)
	Start(id string) (bool, error)
	Stop(id string) (bool, error)
	AssociateSSHKey(id string, sshKeyName string) (bool, error)
	Reboot(id string) (bool, error)
	ChangeComputeOffering(Instance) (bool, error)
	ChangeNetwork(id string, newNetworkId string) (bool, error)
	ResetPassword(id string) (string, error)
	CreateRecoveryPoint(id string, recoveryPoint RecoveryPoint) (bool, error)
}

type InstanceApi struct {
	entityService services.EntityService
}

func NewInstanceService(apiClient api.ApiClient, serviceCode string, environmentName string) InstanceService {
	return &InstanceApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, INSTANCE_ENTITY_TYPE),
	}
}

func parseInstance(data []byte) *Instance {
	instance := Instance{}
	json.Unmarshal(data, &instance)
	return &instance
}

func parseInstanceList(data []byte) []Instance {
	instances := []Instance{}
	json.Unmarshal(data, &instances)
	return instances
}

//Get instance with the specified id for the current environment
func (instanceApi *InstanceApi) Get(id string) (*Instance, error) {
	data, err := instanceApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseInstance(data), nil
}

//List all instances for the current environment
func (instanceApi *InstanceApi) List() ([]Instance, error) {
	return instanceApi.ListWithOptions(map[string]string{})
}

//List all instances for the current environment. Can use options to do sorting and paging.
func (instanceApi *InstanceApi) ListWithOptions(options map[string]string) ([]Instance, error) {
	data, err := instanceApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseInstanceList(data), nil
}

//Create an instance in the current environment
func (instanceApi *InstanceApi) Create(instance Instance) (*Instance, error) {
	send, merr := json.Marshal(instance)
	if merr != nil {
		return nil, merr
	}
	body, err := instanceApi.entityService.Create(send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseInstance(body), nil
}

//Destroy an instance with specified id in the current environment
//Set the purge flag to true if you want to purge immediately
func (instanceApi *InstanceApi) Destroy(id string, purge bool) (bool, error) {
	send, merr := json.Marshal(DestroyOptions{
		PurgeImmediately: purge,
	})
	if merr != nil {
		return false, merr
	}
	_, err := instanceApi.entityService.Delete(id, send, map[string]string{})
	return err == nil, err
}

//Destroy an instance with specified id in the current environment
//Set the purge flag to true if you want to purge immediately
func (instanceApi *InstanceApi) DestroyWithOptions(id string, options DestroyOptions) (bool, error) {
	send, merr := json.Marshal(options)
	if merr != nil {
		return false, merr
	}
	_, err := instanceApi.entityService.Delete(id, send, map[string]string{})
	return err == nil, err
}

//Purge an instance with the specified id in the current environment
//The instance must be in the Destroyed state. To destroy and purge an instance, see the Destroy method
func (instanceApi *InstanceApi) Purge(id string) (bool, error) {
	_, err := instanceApi.entityService.Execute(id, INSTANCE_PURGE_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}

//Recover a destroyed instance with the specified id in the current environment
//Note: Cannot recover instances that have been purged
func (instanceApi *InstanceApi) Recover(id string) (bool, error) {
	_, err := instanceApi.entityService.Execute(id, INSTANCE_RECOVER_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}

//Check if instance with specified id exists in the current environment
func (instanceApi *InstanceApi) Exists(id string) (bool, error) {
	_, err := instanceApi.Get(id)
	if err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok && ccaError.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//Start a stopped instance with specified id exists in the current environment
func (instanceApi *InstanceApi) Start(id string) (bool, error) {
	_, err := instanceApi.entityService.Execute(id, INSTANCE_START_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}

//Stop a running instance with specified id exists in the current environment
func (instanceApi *InstanceApi) Stop(id string) (bool, error) {
	_, err := instanceApi.entityService.Execute(id, INSTANCE_STOP_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}

//Associate an SSH key to the instance with the specified id exists in the current environment
//Note: This will reboot your instance if running
func (instanceApi *InstanceApi) AssociateSSHKey(id string, sshKeyName string) (bool, error) {
	send, merr := json.Marshal(Instance{
		SSHKeyName: sshKeyName,
	})
	if merr != nil {
		return false, merr
	}
	_, err := instanceApi.entityService.Execute(id, INSTANCE_ASSOCIATE_SSH_KEY_OPERATION, send, map[string]string{})
	return err == nil, err
}

//Reboot a running instance with specified id exists in the current environment
func (instanceApi *InstanceApi) Reboot(id string) (bool, error) {
	_, err := instanceApi.entityService.Execute(id, INSTANCE_REBOOT_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}

//Change the compute offering of the instance with the specified id exists in the current environment
//Note: This will reboot your instance if running
func (instanceApi *InstanceApi) ChangeComputeOffering(instance Instance) (bool, error) {
	send, merr := json.Marshal(instance)
	if merr != nil {
		return false, merr
	}
	_, err := instanceApi.entityService.Execute(instance.Id, INSTANCE_CHANGE_COMPUTE_OFFERING_OPERATION, send, map[string]string{})
	return err == nil, err
}

//Reset the password of the instance with the specified id exists in the current environment
func (instanceApi *InstanceApi) ResetPassword(id string) (string, error) {
	body, err := instanceApi.entityService.Execute(id, INSTANCE_RESET_PASSWORD_OPERATION, []byte{}, map[string]string{})
	if err != nil {
		return "", err
	}
	instance := parseInstance(body)
	return instance.Password, nil
}

//Change the network of the instance with the specified id
//Note: This will reboot your instance, remove all pfrs of this instance and remove the instance from all lbrs.
func (instanceApi *InstanceApi) ChangeNetwork(id string, networkId string) (bool, error) {
	send, merr := json.Marshal(Instance{NetworkId: networkId})
	if merr != nil {
		return false, merr
	}
	_, err := instanceApi.entityService.Execute(id, INSTANCE_CHANGE_COMPUTE_OFFERING_OPERATION, send, map[string]string{})
	return err == nil, err
}

//Create a recovery point of the instance with the specified id exists in the current environment
func (instanceApi *InstanceApi) CreateRecoveryPoint(id string, recoveryPoint RecoveryPoint) (bool, error) {
	send, merr := json.Marshal(Instance{
		RecoveryPoint: recoveryPoint,
	})
	if merr != nil {
		return false, merr
	}
	_, err := instanceApi.entityService.Execute(id, INSTANCE_CREATE_RECOVERY_POINT_OPERATION, send, map[string]string{})
	return err == nil, err
}
