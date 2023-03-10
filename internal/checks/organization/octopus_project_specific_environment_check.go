package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"golang.org/x/exp/slices"
	"strings"
)

// OctopusProjectSpecificEnvironmentCheck checks to see if any project variables are unused.
type OctopusProjectSpecificEnvironmentCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusProjectSpecificEnvironmentCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusProjectSpecificEnvironmentCheck {
	return OctopusProjectSpecificEnvironmentCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusProjectSpecificEnvironmentCheck) Id() string {
	return "OctoLintProjectSpecificEnvs"
}

func (o OctopusProjectSpecificEnvironmentCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	projects, err := o.client.Projects.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allLifecycles, err := o.client.Lifecycles.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allEnvironments, err := o.client.Environments.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allChannels, err := o.client.Channels.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	// count the number of times an environment is referenced by a project
	environmentCount := map[string][]string{}
	for _, p := range projects {
		projectEnvironments := []string{}

		// get the allEnvironments from the default lifecycle
		lifecycle := o.getLifecycleById(allLifecycles, p.LifecycleID)
		if lifecycle != nil {
			// this is the default lifecycle, which includes all environments
			if lifecycle.Phases == nil || len(lifecycle.Phases) == 0 {
				projectEnvironments = append(projectEnvironments, o.getAllEnvironmentIds(allEnvironments)...)
			}

			for _, phase := range lifecycle.Phases {
				projectEnvironments = append(projectEnvironments, phase.AutomaticDeploymentTargets...)
				projectEnvironments = append(projectEnvironments, phase.OptionalDeploymentTargets...)
			}
		}

		// get the environments from the channels
		for _, c := range allChannels {
			if c.ProjectID == p.ID {
				channelLifecycle := o.getLifecycleById(allLifecycles, c.LifecycleID)
				if channelLifecycle != nil {
					// this is the default lifecycle, which includes all environments
					if lifecycle.Phases == nil || len(lifecycle.Phases) == 0 {
						projectEnvironments = append(projectEnvironments, o.getAllEnvironmentIds(allEnvironments)...)
					}

					for _, phase := range channelLifecycle.Phases {
						projectEnvironments = append(projectEnvironments, phase.AutomaticDeploymentTargets...)
						projectEnvironments = append(projectEnvironments, phase.OptionalDeploymentTargets...)
					}
				}
			}
		}

		// count the project against the environment
		processedEnvironments := []string{}
		for _, env := range projectEnvironments {
			if slices.Index(processedEnvironments, env) != -1 {
				continue
			}

			if _, ok := environmentCount[env]; !ok {
				environmentCount[env] = []string{}
			}
			environmentCount[env] = append(environmentCount[env], p.Name)
			processedEnvironments = append(processedEnvironments, env)
		}

	}

	// filter down to allEnvironments that have one project
	singleProjectEnvironments := map[string]string{}
	for env, envProjects := range environmentCount {
		if len(envProjects) == 1 {
			singleProjectEnvironments[env] = envProjects[0]
		}
	}

	if len(singleProjectEnvironments) > 0 {
		messages := []string{}
		for env, envProject := range singleProjectEnvironments {
			messages = append(messages, o.getEnvironmentById(allEnvironments, env).Name+" ("+envProject+")")
		}

		return checks.NewOctopusCheckResultImpl(
			"The following environments are used by a single project: "+strings.Join(messages, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no single project environments",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusProjectSpecificEnvironmentCheck) getLifecycleById(lifecycles []*lifecycles.Lifecycle, id string) *lifecycles.Lifecycle {
	for _, l := range lifecycles {
		if l.ID == id {
			return l
		}
	}

	return nil
}

func (o OctopusProjectSpecificEnvironmentCheck) getEnvironmentById(environment []*environments.Environment, id string) *environments.Environment {
	for _, l := range environment {
		if l.ID == id {
			return l
		}
	}

	return nil
}

func (o OctopusProjectSpecificEnvironmentCheck) getAllEnvironmentIds(environment []*environments.Environment) []string {
	ids := []string{}
	for _, l := range environment {
		ids = append(ids, l.ID)
	}

	return ids
}
