package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"strings"
)

const maxStepCount = 20

// OctopusProjectTooManyStepsCheck checks to see if any project has too many steps.
type OctopusProjectTooManyStepsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusProjectTooManyStepsCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusProjectTooManyStepsCheck {
	return OctopusProjectTooManyStepsCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusProjectTooManyStepsCheck) Id() string {
	return "OctoLintTooManySteps"
}

func (o OctopusProjectTooManyStepsCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	projects, err := o.client.Projects.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	complexProjects := []string{}
	for _, p := range projects {
		stepCount, err := o.stepsInDeploymentProcess(p.DeploymentProcessID)

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		if stepCount >= maxStepCount {
			complexProjects = append(complexProjects, p.Name)
		}
	}

	if len(complexProjects) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following projects have 20 or more steps: "+strings.Join(complexProjects, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no projects with too many steps",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusProjectTooManyStepsCheck) stepsInDeploymentProcess(deploymentProcessID string) (int, error) {
	if deploymentProcessID == "" {
		return 0, nil
	}

	resource, err := o.client.DeploymentProcesses.GetByID(deploymentProcessID)

	if err != nil {
		// If we can't find the deployment process, assume zero steps
		if err.(*core.APIError).StatusCode == 404 {
			return 0, nil
		}
		return 0, err
	}

	return len(resource.Steps), nil
}
