package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/octoclient"
)

const maxProjectsInDefaultGroup = 10

type OctopusDefaultProjectGroupCountCheck struct {
	client *client.Client
}

func NewOctopusDefaultProjectGroupCountCheck(client *client.Client) OctopusDefaultProjectGroupCountCheck {
	return OctopusDefaultProjectGroupCountCheck{client: client}
}

func (o OctopusDefaultProjectGroupCountCheck) Id() string {
	return "OctoRecDefaultProjectGroupChildCount"
}

func (o OctopusDefaultProjectGroupCountCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	resource, err := o.client.ProjectGroups.GetByName("Default Project Group")

	if err != nil {
		apiError, ok := err.(*core.APIError)
		if ok && apiError.StatusCode == 404 {
			return checks.NewOctopusCheckResultImpl(
				"The default project group was not found",
				o.Id(),
				"",
				checks.Ok,
				checks.Organization), nil
		}
		return octoclient.ReturnPermissionResultOrError(o.Id(), err)
	}

	if resource != nil {

		projects, err := o.client.ProjectGroups.GetProjects(resource)

		if err != nil {
			return octoclient.ReturnPermissionResultOrError(o.Id(), err)
		}

		if len(projects) > maxProjectsInDefaultGroup {
			return checks.NewOctopusCheckResultImpl(
				"The default project group contains "+fmt.Sprint(len(projects))+" projects. You may want to organize these projects into additional project groups.",
				o.Id(),
				"",
				checks.Warning,
				checks.Organization), nil
		}
	}

	return checks.NewOctopusCheckResultImpl(
		"The number of projects in the default project group is OK",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}
