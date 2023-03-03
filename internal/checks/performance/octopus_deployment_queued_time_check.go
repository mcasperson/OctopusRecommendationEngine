package performance

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"strings"
)

const maxQueueTimeMinutes = 1
const maxQueuedTasks = 10

// OctopusDeploymentQueuedTimeCheck checks to see if any deployments were queued for a long period of time
type OctopusDeploymentQueuedTimeCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusDeploymentQueuedTimeCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusDeploymentQueuedTimeCheck {
	return OctopusDeploymentQueuedTimeCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusDeploymentQueuedTimeCheck) Id() string {
	return "OctoRecDeploymentQueuedTime"
}

func (o OctopusDeploymentQueuedTimeCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	resource, err := o.client.Events.Get(events.EventsQuery{
		EventCategories: []string{"DeploymentQueued", "DeploymentStarted"},
		Skip:            0,
		Take:            1000,
	})

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Performance, err)
	}

	deployments := []string{}
	if resource != nil {
		for _, r := range resource.Items {
			if r.Category == "DeploymentQueued" {
				queuedDeploymentId := o.getDeploymentFromRelatedDocs(r)
				for _, r2 := range resource.Items {
					if r2.Category == "DeploymentStarted" && queuedDeploymentId == o.getDeploymentFromRelatedDocs(r2) {
						queueTime := r2.Occurred.Sub(r.Occurred)
						if queueTime.Minutes() > maxQueueTimeMinutes {
							deployments = append(deployments, queuedDeploymentId)
						}
					}
				}
			}
		}

		if len(deployments) >= maxQueuedTasks {
			return checks.NewOctopusCheckResultImpl(
				fmt.Sprint("Found "+fmt.Sprint(len(deployments)))+" deployments that were queued for longer than "+fmt.Sprintln(maxQueueTimeMinutes)+" minutes. Consider increasing the task cap or adding a HA node to reduce task queue times.",
				o.Id(),
				"",
				checks.Warning,
				checks.Performance), nil
		}
	}

	return checks.NewOctopusCheckResultImpl(
		"Found "+fmt.Sprint(len(deployments))+" deployment tasks that were queued for longer than "+fmt.Sprintln(maxQueueTimeMinutes)+" minutes.",
		o.Id(),
		"",
		checks.Ok,
		checks.Performance), nil
}

func (o OctopusDeploymentQueuedTimeCheck) getDeploymentFromRelatedDocs(event *events.Event) string {
	for _, d := range event.RelatedDocumentIds {
		if strings.HasPrefix(d, "Deployments-") {
			return d
		}
	}
	return ""
}
