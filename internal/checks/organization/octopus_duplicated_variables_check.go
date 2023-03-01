package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	projects2 "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

type projectVar struct {
	project1  *projects2.Project
	variable1 *variables.Variable
	project2  *projects2.Project
	variable2 *variables.Variable
}

type OctopusDuplicatedVariablesCheck struct {
	client *client.Client
}

func NewOctopusDuplicatedVariablesCheck(client *client.Client) OctopusDuplicatedVariablesCheck {
	return OctopusDuplicatedVariablesCheck{client: client}
}

func (o OctopusDuplicatedVariablesCheck) Id() string {
	return "OctoRecDuplicatedVariables"
}

func (o OctopusDuplicatedVariablesCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	projects, err := o.client.Projects.GetAll()

	if err != nil {
		return nil, err
	}

	projectVars := map[*projects2.Project]variables.VariableSet{}
	for _, p := range projects {
		variableSet, err := o.client.Variables.GetAll(p.ID)

		if err != nil {
			return nil, err
		}

		projectVars[p] = variableSet
	}

	duplicateVars := []projectVar{}

	for index1 := 0; index1 < len(projects); index1++ {
		project1 := projects[index1]
		for _, variable1 := range projectVars[project1].Variables {
			if o.shouldIgnoreVariable(variable1) {
				continue
			}

			for index2 := index1 + 1; index2 < len(projects); index2++ {
				project2 := projects[index2]
				for _, variable2 := range projectVars[project2].Variables {
					if variable1.Value == variable2.Value {
						duplicateVars = append(duplicateVars, projectVar{
							project1:  project1,
							variable1: variable1,
							project2:  project2,
							variable2: variable2,
						})
					}
				}
			}
		}
	}

	if len(duplicateVars) > 0 {
		messages := []string{}
		for _, variable := range duplicateVars {
			messages = append(messages, variable.project1.Name+"/"+variable.variable1.Name+" == "+variable.project2.Name+"/"+variable.variable2.Name)
		}

		return checks.NewOctopusCheckResultImpl(
			"The following variables are duplicated between projects. Consider moving these into library variable sets: "+strings.Join(messages, "; "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no duplicated variables",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusDuplicatedVariablesCheck) shouldIgnoreVariable(variable *variables.Variable) bool {
	_, err := strconv.Atoi(variable.Value)
	return variable.Value == "" ||
		variable.Type != "String" ||
		slices.Index(checks.SpecialVars, variable.Name) != -1 ||
		strings.ToLower(variable.Value) == "true" ||
		strings.ToLower(variable.Value) == "false" ||
		err == nil
}
