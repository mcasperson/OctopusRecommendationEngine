package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tasks"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"regexp"
	"strings"
	"time"
)

const maxTimeSinceLastMachineDeployment = time.Hour * 24 * 30

// OctopusUnusedTargetsCheck checks to see if any targets have not been used in a month
type OctopusUnusedTargetsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusUnusedTargetsCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusUnusedTargetsCheck {
	return OctopusUnusedTargetsCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusUnusedTargetsCheck) Id() string {
	return "OctoLintUnusedTargets"
}

func (o OctopusUnusedTargetsCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	targets, err := o.client.Machines.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	unusedMachines := []string{}
	linksTemplate := regexp.MustCompile(`\{.+\}`)
	for _, m := range targets {
		tasksLink := linksTemplate.ReplaceAllString(m.Links["TasksTemplate"], "")
		tasks, err := newclient.Get[resources.Resources[tasks.Task]](o.client.HttpSession(), tasksLink+"?type=Deployment")

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		recentTask := false
		for _, t := range tasks.Items {
			if time.Now().Sub(*t.CompletedTime) < maxTimeSinceLastMachineDeployment {
				recentTask = true
				break
			}
		}

		if !recentTask {
			unusedMachines = append(unusedMachines, m.Name)
		}

	}

	if len(unusedMachines) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following targets have not performed a deployment in 30 days: "+strings.Join(unusedMachines, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no unused targets",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

// naiveStepVariableScan does a simple text search for the variable in a steps properties. This does lead to false positives as simple variables names, like "a",
// will almost certainly appear in a step property text without necessarily being referenced as a variable.
func (o OctopusUnusedTargetsCheck) naiveStepVariableScan(deploymentProcess *deployments.DeploymentProcess, variable *variables.Variable) bool {
	if deploymentProcess != nil {
		for _, s := range deploymentProcess.Steps {
			for _, a := range s.Actions {
				for _, p := range a.Properties {
					if strings.Index(p.Value, variable.Name) != -1 {
						return true
					}
				}
			}
		}
	}

	return false
}

// naiveVariableSetVariableScan does a simple text search for the variable in the value of other variables
func (o OctopusUnusedTargetsCheck) naiveVariableSetVariableScan(variables variables.VariableSet, variable *variables.Variable) bool {
	for _, v := range variables.Variables {
		if strings.Index(v.Value, variable.Name) != -1 {
			return true
		}
	}

	return false
}
