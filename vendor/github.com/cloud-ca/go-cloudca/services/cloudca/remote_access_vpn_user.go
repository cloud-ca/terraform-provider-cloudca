package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

// RemoteAccessVpnUser is an environment wide user for a RemoteAccessVpn
type RemoteAccessVpnUser struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// RemoteAccessVpnUserService defines the interface which is implemented
type RemoteAccessVpnUserService interface {
	Get(id string) (*RemoteAccessVpnUser, error)
	List() ([]RemoteAccessVpnUser, error)
	Create(remoteAccessVpnUser RemoteAccessVpnUser) (bool, error)
	Delete(remoteAccessVpnUser RemoteAccessVpnUser) (bool, error)
}

// RemoteAccessVpnUserApi wraps the EntityService
type RemoteAccessVpnUserApi struct {
	entityService services.EntityService
}

// NewRemoteAccessVpnUserService creates a new VPN User Service for this specific service and environment
func NewRemoteAccessVpnUserService(apiClient api.ApiClient, serviceCode string, environmentName string) RemoteAccessVpnUserService {
	return &RemoteAccessVpnUserApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, REMOTE_ACCESS_VPN_USER_ENTITY_TYPE),
	}
}

func parseRemoteAccessVpnUser(data []byte) *RemoteAccessVpnUser {
	remoteAccessVpnUser := RemoteAccessVpnUser{}
	json.Unmarshal(data, &remoteAccessVpnUser)
	return &remoteAccessVpnUser
}

func parseRemoteAccessVpnUserList(data []byte) []RemoteAccessVpnUser {
	remoteAccessVpnUsers := []RemoteAccessVpnUser{}
	json.Unmarshal(data, &remoteAccessVpnUsers)
	return remoteAccessVpnUsers
}

// Get a specific VPN User in the current environment by their ID
func (remoteAccessVpnUserApi *RemoteAccessVpnUserApi) Get(id string) (*RemoteAccessVpnUser, error) {
	data, err := remoteAccessVpnUserApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseRemoteAccessVpnUser(data), nil
}

// List VPN Users for this environment
func (remoteAccessVpnUserApi *RemoteAccessVpnUserApi) List() ([]RemoteAccessVpnUser, error) {
	data, err := remoteAccessVpnUserApi.entityService.List(map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseRemoteAccessVpnUserList(data), nil
}

// Create a VPN User in the current environment
func (remoteAccessVpnUserApi *RemoteAccessVpnUserApi) Create(remoteAccessVpnUser RemoteAccessVpnUser) (bool, error) {
	send, merr := json.Marshal(remoteAccessVpnUser)
	if merr != nil {
		return false, merr
	}
	_, err := remoteAccessVpnUserApi.entityService.Create(send, map[string]string{})
	return err == nil, err
}

// Delete a specific VPN User in the current environment by their ID
func (remoteAccessVpnUserApi *RemoteAccessVpnUserApi) Delete(remoteAccessVpnUser RemoteAccessVpnUser) (bool, error) {
	send, merr := json.Marshal(remoteAccessVpnUser)
	if merr != nil {
		return false, merr
	}
	_, err := remoteAccessVpnUserApi.entityService.Delete(remoteAccessVpnUser.Id, send, map[string]string{})
	return err == nil, err
}
