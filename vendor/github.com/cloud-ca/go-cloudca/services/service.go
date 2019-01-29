package services

//Holds all the resources available for a particular service
type ServiceResources interface {
	GetServiceType() string
}
