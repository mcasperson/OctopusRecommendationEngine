package factory

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/organization"
)

type OctopusCheckFactory struct {
	client *client.Client
}

func NewOctopusCheckFactory(client *client.Client) OctopusCheckFactory {
	return OctopusCheckFactory{client: client}
}

// BuildAllChecks creates new instances of all the checks and returns them as an array.
func (o OctopusCheckFactory) BuildAllChecks() ([]checks.OctopusCheck, error) {
	return []checks.OctopusCheck{
		organization.NewOctopusEnvironmentCountCheck(o.client),
		organization.NewOctopusDefaultProjectGroupCountCheck(o.client),
		organization.NewOctopusEmptyProjectCheck(o.client),
	}, nil
}
