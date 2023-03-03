package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"strings"
)

type OctopusLifecycleRetentionPolicyCheck struct {
	client *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusLifecycleRetentionPolicyCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusLifecycleRetentionPolicyCheck {
	return OctopusLifecycleRetentionPolicyCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusLifecycleRetentionPolicyCheck) Id() string {
	return "OctoRecLifecycleRetention"
}

func (o OctopusLifecycleRetentionPolicyCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	lifecycles, err := o.client.Lifecycles.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	keepsForever := []string{}
	for _, l := range lifecycles {
		phaseKeepsForever, err := o.anyPhasesKeepForever(l.Phases)

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		lifecycleKeepsForever := l.ReleaseRetentionPolicy.ShouldKeepForever || l.TentacleRetentionPolicy.ShouldKeepForever

		if lifecycleKeepsForever || phaseKeepsForever {
			keepsForever = append(keepsForever, l.Name)
		}
	}

	if len(keepsForever) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following lifecycles have retention policies that keep releases or files forever: "+strings.Join(keepsForever, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no lifecycles with retention policies that keep releases or files forever",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusLifecycleRetentionPolicyCheck) anyPhasesKeepForever(phases []*lifecycles.Phase) (bool, error) {
	if len(phases) == 0 {
		return false, nil
	}

	for _, p := range phases {
		keepReleasesForver := p.ReleaseRetentionPolicy != nil && p.ReleaseRetentionPolicy.ShouldKeepForever
		keepFilesForever := p.TentacleRetentionPolicy != nil && p.TentacleRetentionPolicy.ShouldKeepForever

		if keepReleasesForver || keepFilesForever {
			return true, nil
		}
	}

	return false, nil
}
