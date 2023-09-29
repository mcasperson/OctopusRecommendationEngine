package security

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"strings"
)

// CustomProject is the simplest representation of a project and its version controlled settings
type CustomProject struct {
	PersistenceSettings CustomPersistenceSettings `json:"PersistenceSettings"`
	Name                string                    `json:"Name"`
}

type CustomPersistenceSettings struct {
	Type        string             `json:"Type"`
	Credentials *CustomCredentials `json:"Credentials"`
}

type CustomCredentials struct {
	Type     string  `json:"Type"`
	Username *string `json:"Username"`
}

// OctopusDuplicatedGitCredentialsCheck reports on any perpetual api keys
type OctopusDuplicatedGitCredentialsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusDuplicatedGitCredentialsCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusDuplicatedGitCredentialsCheck {
	return OctopusDuplicatedGitCredentialsCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusDuplicatedGitCredentialsCheck) Id() string {
	return "OctoLintSharedGitUsername"
}

func (o OctopusDuplicatedGitCredentialsCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	o.client.Users.GetAll()

	allProjects, err := newclient.Get[resources.Resources[CustomProject]](o.client.HttpSession(), o.client.HttpSession().BaseURL.String()+"/Projects?take=1000")

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	gitUsernameCounts := map[string]int{}
	gitUsernameProjects := map[string][]string{}
	for _, p := range allProjects.Items {
		if p.PersistenceSettings.Type == "VersionControlled" &&
			p.PersistenceSettings.Credentials.Type == "UsernamePassword" &&
			p.PersistenceSettings.Credentials.Username != nil {
			if _, ok := gitUsernameCounts[*p.PersistenceSettings.Credentials.Username]; !ok {
				gitUsernameCounts[*p.PersistenceSettings.Credentials.Username] = 0
			}
			gitUsernameCounts[*p.PersistenceSettings.Credentials.Username]++

			if _, ok := gitUsernameProjects[*p.PersistenceSettings.Credentials.Username]; !ok {
				gitUsernameProjects[*p.PersistenceSettings.Credentials.Username] = []string{}
			}
			gitUsernameProjects[*p.PersistenceSettings.Credentials.Username] = append(gitUsernameProjects[*p.PersistenceSettings.Credentials.Username], p.Name)
		}
	}

	duplicatedGitCredentials := map[string][]string{}
	for u, c := range gitUsernameCounts {
		if c > 1 {
			duplicatedGitCredentials[u] = gitUsernameProjects[u]
		}
	}

	if len(duplicatedGitCredentials) != 0 {
		message := []string{}
		for u, p := range duplicatedGitCredentials {
			message = append(message, u+" ("+strings.Join(p, ", ")+")")
		}

		return checks.NewOctopusCheckResultImpl(
			"The following Git usernames have been reused across the following projects:\n"+strings.Join(message, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"No Git usernames have been resued",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}
