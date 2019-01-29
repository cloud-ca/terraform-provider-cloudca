package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type SSHKey struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PublicKey   string `json:"publicKey,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

type SSHKeyService interface {
	Get(name string) (*SSHKey, error)
	List() ([]SSHKey, error)
	ListWithOptions(options map[string]string) ([]SSHKey, error)
	Create(key SSHKey) (*SSHKey, error)
	Delete(id string) (bool, error)
}

type SSHKeyApi struct {
	entityService services.EntityService
}

func NewSSHKeyService(apiClient api.ApiClient, serviceCode string, environmentName string) SSHKeyService {
	return &SSHKeyApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, SSH_KEY_ENTITY_TYPE),
	}
}

func parseSSHKey(data []byte) *SSHKey {
	var sshKey SSHKey
	json.Unmarshal(data, &sshKey)
	return &sshKey
}

func parseSSHKeyList(data []byte) []SSHKey {
	sshKeys := []SSHKey{}
	json.Unmarshal(data, &sshKeys)
	return sshKeys
}

//Get SSH key with the specified id for the current environment
func (sshKeyApi *SSHKeyApi) Get(name string) (*SSHKey, error) {
	data, err := sshKeyApi.entityService.Get(name, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseSSHKey(data), nil
}

//List all SSH keys for the current environment
func (sshKeyApi *SSHKeyApi) List() ([]SSHKey, error) {
	return sshKeyApi.ListWithOptions(map[string]string{})
}

//List all SSH keys for the current environment. Can use options to do sorting and paging.
func (sshKeyApi *SSHKeyApi) ListWithOptions(options map[string]string) ([]SSHKey, error) {
	data, err := sshKeyApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseSSHKeyList(data), nil
}

//Create an SSH key in the current environment
func (sshKeyApi *SSHKeyApi) Create(key SSHKey) (*SSHKey, error) {
	send, merr := json.Marshal(key)
	if merr != nil {
		return nil, merr
	}
	body, err := sshKeyApi.entityService.Create(send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseSSHKey(body), nil
}

//Delete an SSH Key with specified id in the current environment
func (sshKeyApi *SSHKeyApi) Delete(id string) (bool, error) {
	_, err := sshKeyApi.entityService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}
