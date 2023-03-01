package factory

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/organization"
)

type OctopusCheckFactory struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusCheckFactory(client *client.Client) OctopusCheckFactory {
	return OctopusCheckFactory{client: client, errorHandler: checks.OctopusClientPermissiveErrorHandler{}}
}

// BuildAllChecks creates new instances of all the checks and returns them as an array.
func (o OctopusCheckFactory) BuildAllChecks() ([]checks.OctopusCheck, error) {
	return []checks.OctopusCheck{
		organization.NewOctopusEnvironmentCountCheck(o.client, o.errorHandler),
		organization.NewOctopusDefaultProjectGroupCountCheck(o.client, o.errorHandler),
		organization.NewOctopusEmptyProjectCheck(o.client, o.errorHandler),
		organization.NewOctopusUnusedVariablesCheck(o.client, o.errorHandler),
		organization.NewOctopusDuplicatedVariablesCheck(o.client, o.errorHandler),
		organization.NewOctopusProjectTooManyStepsCheck(o.client, o.errorHandler),
	}, nil
}
