package factory

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/organization"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/performance"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/security"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"strings"
)

// OctopusCheckFactory builds all the lint checks. This is where you can customize things like error handlers.
type OctopusCheckFactory struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	url          string
	space        string
}

func NewOctopusCheckFactory(client *client.Client, url string, space string) OctopusCheckFactory {
	return OctopusCheckFactory{client: client, url: url, space: space, errorHandler: checks.OctopusClientPermissiveErrorHandler{}}
}

// BuildAllChecks creates new instances of all the checks and returns them as an array.
func (o OctopusCheckFactory) BuildAllChecks(skipChecks string) ([]checks.OctopusCheck, error) {
	skipChecksSlice := lo.Map(strings.Split(skipChecks, ","), func(item string, index int) string {
		return strings.TrimSpace(item)
	})

	allChecks := []checks.OctopusCheck{
		organization.NewOctopusEnvironmentCountCheck(o.client, o.errorHandler),
		organization.NewOctopusDefaultProjectGroupCountCheck(o.client, o.errorHandler),
		organization.NewOctopusEmptyProjectCheck(o.client, o.errorHandler),
		organization.NewOctopusUnusedVariablesCheck(o.client, o.errorHandler),
		organization.NewOctopusDuplicatedVariablesCheck(o.client, o.errorHandler),
		organization.NewOctopusProjectTooManyStepsCheck(o.client, o.errorHandler),
		organization.NewOctopusLifecycleRetentionPolicyCheck(o.client, o.errorHandler),
		organization.NewOctopusUnusedTargetsCheck(o.client, o.errorHandler),
		organization.NewOctopusProjectSpecificEnvironmentCheck(o.client, o.errorHandler),
		organization.NewOctopusTenantsInsteadOfTagsCheck(o.client, o.errorHandler),
		organization.NewOctopusProjectGroupsWithExclusiveEnvironmentsCheck(o.client, o.errorHandler),
		organization.NewOctopusUnhealthyTargetCheck(o.client, o.errorHandler),
		security.NewOctopusDeploymentQueuedByAdminCheck(o.client, o.errorHandler),
		security.NewOctopusPerpetualApiKeysCheck(o.client, o.errorHandler),
		security.NewOctopusDuplicatedGitCredentialsCheck(o.client, o.errorHandler),
		performance.NewOctopusDeploymentQueuedTimeCheck(o.client, o.url, o.space, o.errorHandler),
	}

	return lo.Filter(allChecks, func(item checks.OctopusCheck, index int) bool {
		return slices.Index(skipChecksSlice, item.Id()) == -1
	}), nil
}
