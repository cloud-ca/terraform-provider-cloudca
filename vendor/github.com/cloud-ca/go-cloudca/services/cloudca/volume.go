package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

const (
	VOLUME_TYPE_OS   = "OS"
	VOLUME_TYPE_DATA = "DATA"
)

type Volume struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
	CreationDate     string `json:"creationDate,omitempty"`
	Size             int    `json:"size,omitempty"`
	GbSize           int    `json:"sizeInGb,omitempty"`
	DiskOfferingId   string `json:"diskOfferingId,omitempty"`
	DiskOfferingName string `json:"diskOfferingName,omitempty"`
	TemplateId       string `json:"templateId,omitempty"`
	ZoneName         string `json:"zoneName,omitempty"`
	ZoneId           string `json:"zoneId,omitempty"`
	State            string `json:"state,omitempty"`
	InstanceName     string `json:"instanceName,omitempty"`
	InstanceId       string `json:"instanceId,omitempty"`
	InstanceState    string `json:"instanceState,omitempty"`
	Iops             int    `json:"iops,omitempty"`
}

type VolumeService interface {
	Get(id string) (*Volume, error)
	List() ([]Volume, error)
	ListOfType(volumeType string) ([]Volume, error)
	ListWithOptions(options map[string]string) ([]Volume, error)
	Create(Volume) (*Volume, error)
	Resize(*Volume) error
	Delete(string) error
	AttachToInstance(*Volume, string) error
	DetachFromInstance(*Volume) error
}

type VolumeApi struct {
	entityService services.EntityService
}

func NewVolumeService(apiClient api.ApiClient, serviceCode string, environmentName string) VolumeService {
	return &VolumeApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, VOLUME_ENTITY_TYPE),
	}
}

func parseVolume(data []byte) *Volume {
	volume := Volume{}
	json.Unmarshal(data, &volume)
	return &volume
}

func parseVolumeList(data []byte) []Volume {
	volumes := []Volume{}
	json.Unmarshal(data, &volumes)
	return volumes
}

//Get volume with the specified id for the current environment
func (volumeApi *VolumeApi) Get(id string) (*Volume, error) {
	data, err := volumeApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseVolume(data), nil
}

//List all volumes for the current environment
func (volumeApi *VolumeApi) List() ([]Volume, error) {
	return volumeApi.ListWithOptions(map[string]string{})
}

//List all volumes of specified type for the current environment
func (volumeApi *VolumeApi) ListOfType(volumeType string) ([]Volume, error) {
	return volumeApi.ListWithOptions(map[string]string{
		"type": volumeType,
	})
}

//List all volumes for the current environment. Can use options to do sorting and paging.
func (volumeApi *VolumeApi) ListWithOptions(options map[string]string) ([]Volume, error) {
	data, err := volumeApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseVolumeList(data), nil
}

func (api *VolumeApi) Create(volume Volume) (*Volume, error) {
	body, err := json.Marshal(volume)
	if err != nil {
		return nil, err
	}
	res, err := api.entityService.Create(body, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseVolume(res), nil
}

func (api *VolumeApi) Delete(volumeId string) error {
	_, err := api.entityService.Delete(volumeId, []byte{}, map[string]string{})
	return err
}

func (api *VolumeApi) Resize(volume *Volume) error {
	body, err := json.Marshal(volume)
	if err != nil {
		return err
	}
	_, err = api.entityService.Execute(volume.Id, "resize", body, map[string]string{})
	return err
}

func (api *VolumeApi) AttachToInstance(volume *Volume, instanceId string) error {
	body, err := json.Marshal(Volume{
		InstanceId: instanceId,
	})
	if err != nil {
		return err
	}
	_, err = api.entityService.Execute(volume.Id, "attachToInstance", body, map[string]string{})
	return err
}

func (api *VolumeApi) DetachFromInstance(volume *Volume) error {
	_, err := api.entityService.Execute(volume.Id, "detachFromInstance", []byte{}, map[string]string{})
	return err
}
