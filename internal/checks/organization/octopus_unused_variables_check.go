package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	projects2 "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"golang.org/x/exp/slices"
	"strings"
)

// OctopusUnusedVariablesCheck checks to see if any project variables are unused.
type OctopusUnusedVariablesCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusUnusedVariablesCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusUnusedVariablesCheck {
	return OctopusUnusedVariablesCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusUnusedVariablesCheck) Id() string {
	return "OctoLintUnusedVariables"
}

func (o OctopusUnusedVariablesCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	projects, err := o.client.Projects.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	unusedVars := map[*projects2.Project][]*variables.Variable{}
	for _, p := range projects {
		deploymentProcess, err := o.client.DeploymentProcesses.GetByID(p.DeploymentProcessID)

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		variableSet, err := o.client.Variables.GetAll(p.ID)

		if err != nil {
			return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
		}

		for _, v := range variableSet.Variables {
			if slices.Index(checks.SpecialVars[:], v.Name) != -1 {
				continue
			}

			if !(o.naiveStepVariableScan(deploymentProcess, v) || o.naiveVariableSetVariableScan(variableSet, v)) {
				if _, ok := unusedVars[p]; !ok {
					unusedVars[p] = []*variables.Variable{}
				}
				unusedVars[p] = append(unusedVars[p], v)
			}
		}
	}

	if len(unusedVars) > 0 {
		messages := []string{}
		for p, variables := range unusedVars {
			if len(variables) != 0 {
				for _, variable := range variables {
					messages = append(messages, p.Name+"/"+variable.Name)
				}
			}
		}

		return checks.NewOctopusCheckResultImpl(
			"The following variables are unused: "+strings.Join(messages, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no unused variables",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

// naiveStepVariableScan does a simple text search for the variable in a steps properties. This does lead to false positives as simple variables names, like "a",
// will almost certainly appear in a step property text without necessarily being referenced as a variable.
func (o OctopusUnusedVariablesCheck) naiveStepVariableScan(deploymentProcess *deployments.DeploymentProcess, variable *variables.Variable) bool {
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
func (o OctopusUnusedVariablesCheck) naiveVariableSetVariableScan(variables variables.VariableSet, variable *variables.Variable) bool {
	for _, v := range variables.Variables {
		if strings.Index(v.Value, variable.Name) != -1 {
			return true
		}
	}

	return false
}
