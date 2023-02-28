package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
)

const maxEnvironments = 20

type OctopusEnvironmentCountCheck struct {
	client *client.Client
}

func NewOctopusEnvironmentCountCheck(client *client.Client) OctopusEnvironmentCountCheck {
	return OctopusEnvironmentCountCheck{client: client}
}

func (o OctopusEnvironmentCountCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	query := environments.EnvironmentsQuery{
		PartialName: "",
		Skip:        0,
		Take:        1000,
	}
	resources, err := o.client.Environments.Get(query)

	if err != nil {
		return nil, err
	}

	if len(resources.Items) > maxEnvironments {
		return checks.NewOctopusCheckResultImpl(
			"The recommended maximum number of environments is "+fmt.Sprint(maxEnvironments)+". You have at least "+fmt.Sprint(len(resources.Items)),
			"OctoRecEnvironmentCount",
			"https://octopus.com/docs/getting-started/best-practices/environments-and-deployment-targets-and-roles#environments",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"The number of environments in the space is OK",
		"OctoRecEnvironmentCount",
		"https://octopus.com/docs/getting-started/best-practices/environments-and-deployment-targets-and-roles#environments",
		checks.Ok,
		checks.Organization), nil
}
