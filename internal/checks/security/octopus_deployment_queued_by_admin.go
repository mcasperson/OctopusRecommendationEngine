package security

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"golang.org/x/exp/slices"
	"strings"
)

type OctopusDeploymentQueuedByAdminCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusDeploymentQueuedByAdminCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusDeploymentQueuedByAdminCheck {
	return OctopusDeploymentQueuedByAdminCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusDeploymentQueuedByAdminCheck) Id() string {
	return "OctoRecDeploymentQueuedByAdmin"
}

func (o OctopusDeploymentQueuedByAdminCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	resource, err := o.client.Events.Get(events.EventsQuery{
		EventCategories: []string{"DeploymentQueued"},
		Skip:            0,
		Take:            1000,
	})

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	if resource != nil {
		projectsDeployedByAdmins := []string{}
		projects := []string{}
		for _, r := range resource.Items {
			projectId := o.getProjectFromRelatedDocs(r)

			if projectId == "" {
				continue
			}

			project, err := o.client.Projects.GetByID(projectId)

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					return nil, err
				}
				continue
			}

			if slices.Index(projects, project.ID) != -1 {
				continue
			}

			projects = append(projects, project.ID)

			user, err := o.client.Users.Get(users.UsersQuery{
				Filter: r.Username,
				Skip:   0,
				Take:   1,
			})

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					return nil, err
				}
				continue
			}

			usersWhoDeployedProject := []string{}
			for _, u := range user.Items {
				teams, err := o.getAdminTeams()

				if err != nil {
					if !o.errorHandler.ShouldContinue(err) {
						return nil, err
					}
					continue
				}

				for _, t := range teams {
					if slices.Index(t.MemberUserIDs, u.ID) != -1 && slices.Index(usersWhoDeployedProject, u.Username) == -1 {
						usersWhoDeployedProject = append(usersWhoDeployedProject, u.Username)
					}
				}
			}

			result := project.Name + "(" + strings.Join(usersWhoDeployedProject, ",") + ")"
			if slices.Index(projectsDeployedByAdmins, result) == -1 {
				projectsDeployedByAdmins = append(projectsDeployedByAdmins, project.Name+" ("+strings.Join(usersWhoDeployedProject, ",")+")")
			}
		}

		if len(projectsDeployedByAdmins) != 0 {
			return checks.NewOctopusCheckResultImpl(
				"The following projects were deployed by admins. Consider creating a limited user account to perform deployments: "+strings.Join(projectsDeployedByAdmins, ", "),
				o.Id(),
				"",
				checks.Warning,
				checks.Security), nil
		}
	}

	return checks.NewOctopusCheckResultImpl(
		"No deployments were found",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}

func (o OctopusDeploymentQueuedByAdminCheck) getProjectFromRelatedDocs(event *events.Event) string {
	for _, d := range event.RelatedDocumentIds {
		if strings.HasPrefix(d, "Projects-") {
			return d
		}
	}
	return ""
}

func (o OctopusDeploymentQueuedByAdminCheck) getAdminTeams() ([]*teams.Team, error) {
	adminTeams := []string{"Octopus Administrators", "Space Managers", "Octopus Managers"}

	teamResources := []*teams.Team{}
	for _, adminTeam := range adminTeams {
		team, err := o.client.Teams.Get(teams.TeamsQuery{
			IDs:           nil,
			IncludeSystem: true,
			PartialName:   adminTeam,
			Skip:          0,
			Take:          1,
		})

		if err != nil {
			return nil, err
		}

		teamResources = append(teamResources, team.Items...)
	}

	return teamResources, nil
}
