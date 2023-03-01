package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/octoclient"
	"strings"
)

type OctopusEmptyProjectCheck struct {
	client *client.Client
}

func NewOctopusEmptyProjectCheck(client *client.Client) OctopusEmptyProjectCheck {
	return OctopusEmptyProjectCheck{client: client}
}

func (o OctopusEmptyProjectCheck) Id() string {
	return "OctoRecEmptyProject"
}

func (o OctopusEmptyProjectCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	projects, err := o.client.Projects.GetAll()

	if err != nil {
		return octoclient.ReturnPermissionResultOrError(o.Id(), err)
	}

	runbooks, err := o.client.Runbooks.GetAll()

	emptyProjects := []string{}
	for _, p := range projects {
		stepCount, err := o.stepsInDeploymentProcess(p.DeploymentProcessID)

		if err != nil {
			return octoclient.ReturnPermissionResultOrError(o.Id(), err)
		}

		if runbooksInProject(p.ID, runbooks) == 0 && stepCount == 0 {
			emptyProjects = append(emptyProjects, p.Name)
		}
	}

	if len(emptyProjects) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following projects have no runbooks and no deployment process: "+strings.Join(emptyProjects, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no empty projects",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func runbooksInProject(projectID string, runbooks []*runbooks.Runbook) int {
	count := 0
	for _, r := range runbooks {
		if r.ProjectID == projectID {
			count++
		}
	}
	return count
}

func (o OctopusEmptyProjectCheck) stepsInDeploymentProcess(deploymentProcessID string) (int, error) {
	if deploymentProcessID == "" {
		return 0, nil
	}

	resource, err := o.client.DeploymentProcesses.GetByID(deploymentProcessID)

	if err != nil {
		// If we can't find the deployment process, assume zero steps
		apiError, ok := err.(*core.APIError)
		if ok && apiError.StatusCode == 404 {
			return 0, nil
		}
		return 0, err
	}

	return len(resource.Steps), nil
}
