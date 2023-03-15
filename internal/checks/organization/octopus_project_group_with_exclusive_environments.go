package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"golang.org/x/exp/slices"
	"strings"
)

// OctopusProjectGroupsWithExclusiveEnvironmentsCheck checks to see if the project groups contain projects that have mutually exclusive environments.
type OctopusProjectGroupsWithExclusiveEnvironmentsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusProjectGroupsWithExclusiveEnvironmentsCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusProjectGroupsWithExclusiveEnvironmentsCheck {
	return OctopusProjectGroupsWithExclusiveEnvironmentsCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) Id() string {
	return "OctoLintProjectGroupsWithExclusiveEnvironments"
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	allProjectGroups, err := o.client.ProjectGroups.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allProjects, err := o.client.Projects.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allLifecycles, err := o.client.Lifecycles.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	projectGroupsWithExclusiveEnvs := []string{}
	for _, pg := range allProjectGroups {
		// Find the groups of environments captured in the default lifecyles of the projects in the project group
		envGroups := [][]string{}
		for _, p := range allProjects {
			if p.ProjectGroupID == pg.ID {
				lifecycle := o.getLifecycleById(allLifecycles, p.LifecycleID)
				projectEnvironments := o.getLifecycleEnvironments(lifecycle)
				envGroups = append(envGroups, projectEnvironments)
			}
		}

		// don't do any further processing if there was just one project
		if len(envGroups) <= 1 {
			continue
		}

		// Attempt to find at least every environment in a lifecycle with one in another environment
		for i, eg1 := range envGroups[0 : len(envGroups)-1] {
			allExclusive := true

			for _, eg2 := range envGroups[i+1:] {
				for _, e1 := range eg1 {
					if slices.Index(eg2, e1) != -1 {
						allExclusive = false
						break
					}
				}
			}

			// if none of the environments from this lifecycle are found in any other lifecycles, we have an project with exclusive environments
			if allExclusive && slices.Index(projectGroupsWithExclusiveEnvs, pg.Name) == -1 {
				projectGroupsWithExclusiveEnvs = append(projectGroupsWithExclusiveEnvs, pg.Name)
			}
		}
	}

	if len(projectGroupsWithExclusiveEnvs) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following project groups contain projects with mutually exclusive environments in their default lifecycle: "+strings.Join(projectGroupsWithExclusiveEnvs, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no project groups with mutually exclusive lifecycles",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) getLifecycleEnvironments(lifecycle *lifecycles.Lifecycle) []string {
	projectEnvironments := []string{}
	for _, phase := range lifecycle.Phases {
		projectEnvironments = append(projectEnvironments, phase.AutomaticDeploymentTargets...)
		projectEnvironments = append(projectEnvironments, phase.OptionalDeploymentTargets...)
	}
	slices.Sort(projectEnvironments)
	return projectEnvironments
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) getLifecycleById(lifecycles []*lifecycles.Lifecycle, id string) *lifecycles.Lifecycle {
	for _, l := range lifecycles {
		if l.ID == id {
			return l
		}
	}

	return nil
}
