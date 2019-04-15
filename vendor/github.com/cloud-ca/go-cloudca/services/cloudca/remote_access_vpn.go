package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

const (
	REMOTE_ACCESS_VPN_ENABLE_OPERATION  = "enable"
	REMOTE_ACCESS_VPN_DISABLE_OPERATION = "disable"
)

// RemoteAccessVpn is a VPN configuration which can connect a RemoteAccessVpnUser
type RemoteAccessVpn struct {
	Certificate       string `json:"certificate,omitempty"`
	Id                string `json:"id,omitempty"`
	PresharedKey      string `json:"presharedKey,omitempty"`
	PublicIpAddress   string `json:"publicIpAddress,omitempty"`
	PublicIpAddressId string `json:"publicIpAddressId,omitempty"`
	State             string `json:"state,omitempty"`
	Type              string `json:"type,omitempty"`
}

// RemoteAccessVpnService defines the interface which is implemented
type RemoteAccessVpnService interface {
	Get(id string) (*RemoteAccessVpn, error)
	List() ([]RemoteAccessVpn, error)
	ListWithOptions(options map[string]string) ([]RemoteAccessVpn, error)
	Enable(id string) (bool, error)
	Disable(id string) (bool, error)
}

// RemoteAccessVpnApi wraps the EntityService
type RemoteAccessVpnApi struct {
	entityService services.EntityService
}

// NewRemoteAccessVpnService creates a new VPN Service for this specific service and environment
func NewRemoteAccessVpnService(apiClient api.ApiClient, serviceCode string, environmentName string) RemoteAccessVpnService {
	return &RemoteAccessVpnApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, REMOTE_ACCESS_VPN_ENTITY_TYPE),
	}
}

func parseRemoteAccessVpn(data []byte) *RemoteAccessVpn {
	remoteAccessVpn := RemoteAccessVpn{}
	json.Unmarshal(data, &remoteAccessVpn)
	return &remoteAccessVpn
}

func parseRemoteAccessVpnList(data []byte) []RemoteAccessVpn {
	remoteAccessVpns := []RemoteAccessVpn{}
	json.Unmarshal(data, &remoteAccessVpns)
	return remoteAccessVpns
}

// Get a specific VPN in the current environment by its ID
func (remoteAccessVpnApi *RemoteAccessVpnApi) Get(id string) (*RemoteAccessVpn, error) {
	data, err := remoteAccessVpnApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseRemoteAccessVpn(data), nil
}

// List the available VPNs in the current environment
func (remoteAccessVpnApi *RemoteAccessVpnApi) List() ([]RemoteAccessVpn, error) {
	return remoteAccessVpnApi.ListWithOptions(map[string]string{})
}

// ListWithOptions lists the available VPNs in the current environment with options
func (remoteAccessVpnApi *RemoteAccessVpnApi) ListWithOptions(options map[string]string) ([]RemoteAccessVpn, error) {
	data, err := remoteAccessVpnApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseRemoteAccessVpnList(data), nil
}

// Enable a specific VPN in the current environment by its ID
func (remoteAccessVpnApi *RemoteAccessVpnApi) Enable(id string) (bool, error) {
	_, err := remoteAccessVpnApi.entityService.Execute(id, REMOTE_ACCESS_VPN_ENABLE_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}

// Disable a specific VPN in the current environment by its ID
func (remoteAccessVpnApi *RemoteAccessVpnApi) Disable(id string) (bool, error) {
	_, err := remoteAccessVpnApi.entityService.Execute(id, REMOTE_ACCESS_VPN_DISABLE_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}
