package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"strings"
	"time"
)

const maxHealthCheckTime = time.Hour * 24 * 30

// OctopusUnhealthyTargetCheck find targets that have not been healthy in the last 30 days.
type OctopusUnhealthyTargetCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusUnhealthyTargetCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusUnhealthyTargetCheck {
	return OctopusUnhealthyTargetCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusUnhealthyTargetCheck) Id() string {
	return "OctoLintUnhealthyTargets"
}

func (o OctopusUnhealthyTargetCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	allMachines, err := o.client.Machines.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	unhealthyMachines := []string{}
	for _, m := range allMachines {
		wasEverHealthy := true
		if m.HealthStatus == "Unhealthy" {
			wasEverHealthy = false

			targetEvents, err := o.client.Events.Get(events.EventsQuery{
				Regarding: m.ID,
			})

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					return nil, err
				}
				continue
			}

			for _, e := range targetEvents.Items {
				if e.Category == "MachineHealthy" && time.Now().Sub(e.Occurred) < maxHealthCheckTime {
					wasEverHealthy = true
					break
				}
			}
		}

		if !wasEverHealthy {
			unhealthyMachines = append(unhealthyMachines, m.Name)
		}
	}

	if len(unhealthyMachines) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following targets have not been healthy in the last 30 days:\n"+strings.Join(unhealthyMachines, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no targets that were unhealthy for all of the last 30 days",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}
