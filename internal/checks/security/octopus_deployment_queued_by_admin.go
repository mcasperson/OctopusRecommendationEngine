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
	"time"
)

// OctopusDeploymentQueuedByAdminCheck checks to see if any deployments were initiated by someone from the admin teams.
// This usually means that a more specific and limited user should be created to perform deployments.
type OctopusDeploymentQueuedByAdminCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusDeploymentQueuedByAdminCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusDeploymentQueuedByAdminCheck {
	return OctopusDeploymentQueuedByAdminCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusDeploymentQueuedByAdminCheck) Id() string {
	return "OctoLintDeploymentQueuedByAdmin"
}

func (o OctopusDeploymentQueuedByAdminCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	projects, err := o.client.Projects.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	teams, err := o.getAdminTeams()
	
	if err != nil {
		if !o.errorHandler.ShouldContinue(err) {
			return nil, err
		}
	}

	projectsDeployedByAdmins := []string{}

	now := time.Now()
	fromDate := now.AddDate(0, -3, 0)
	from := fromDate.Format("2006-01-02")

	for _, p := range projects {
		projectId := p.ID
		usersWhoDeployedProject := []string{}

		resource, err := o.client.Events.Get(events.EventsQuery{
			EventCategories: []string{"DeploymentQueued"},
			Projects: []string{projectId},
			Skip:            0,
			Take:            100,
			From:			 from,
		})
	
		if err != nil {
			return o.errorHandler.HandleError(o.Id(), checks.Security, err)
		}
	
		if resource != nil {
			for _, r := range resource.Items {
				if r.Username == "system" {
					continue
				}

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
	
				for _, u := range user.Items {
					for _, t := range teams {
						if slices.Index(t.MemberUserIDs, u.ID) != -1 && slices.Index(usersWhoDeployedProject, u.Username) == -1 {
							usersWhoDeployedProject = append(usersWhoDeployedProject, u.Username)
						}
					}
				}
			}
		}

		if len(usersWhoDeployedProject) != 0 {
			projectsDeployedByAdmins = append(projectsDeployedByAdmins, p.Name+" ("+strings.Join(usersWhoDeployedProject, ",")+")")
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

	return checks.NewOctopusCheckResultImpl(
		"No deployments were found",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
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
