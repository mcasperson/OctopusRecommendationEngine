package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"golang.org/x/exp/slices"
	"strings"
)

// OctopusTenantsInsteadOfTagsCheck checks to see if any common groups of tenants are found against common resources like accounts, targets etc
type OctopusTenantsInsteadOfTagsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusTenantsInsteadOfTagsCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusTenantsInsteadOfTagsCheck {
	return OctopusTenantsInsteadOfTagsCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusTenantsInsteadOfTagsCheck) Id() string {
	return "OctoRecDirectTenantReferences"
}

func (o OctopusTenantsInsteadOfTagsCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	allTenants, err := o.client.Tenants.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allAccounts, err := o.client.Accounts.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	tenantReferences := map[string]int{}
	for _, a := range allAccounts {
		if a.GetTenantedDeploymentMode() == core.TenantedDeploymentModeTenantedOrUntenanted {
			tenantIds := a.GetTenantIDs()
			slices.Sort(tenantIds)
			tenants := strings.Join(tenantIds, ",")

			if _, ok := tenantReferences[tenants]; !ok {
				tenantReferences[tenants] = 0
			}
			tenantReferences[tenants]++
		}
	}

	// get any commonly grouped tenants
	multipleTenantReferences := []string{}
	for tenantGroups, groupCount := range tenantReferences {
		if groupCount > 1 {
			multipleTenantReferences = append(multipleTenantReferences, tenantGroups)
		}
	}

	if len(multipleTenantReferences) > 0 {

		// We have to convert the comma separated list of tenant IDs into a comma separated list of tenant names
		groupedTenants := []string{}
		for _, groupedTenant := range multipleTenantReferences {
			splitTenants := strings.Split(groupedTenant, ",")
			splitTenantNames := []string{}
			for _, splitTenant := range splitTenants {
				splitTenantNames = append(splitTenantNames, o.getTenantNameById(allTenants, splitTenant))
			}
			groupedTenants = append(groupedTenants, strings.Join(splitTenantNames, ", "))
		}

		return checks.NewOctopusCheckResultImpl(
			"The following groups of tenants have been directly referenced more than once, and may be better grouped as tenant tags: "+strings.Join(groupedTenants, "; "),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"No duplicate groups of tenants were found",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusTenantsInsteadOfTagsCheck) getTenantNameById(tenants []*tenants.Tenant, id string) string {
	for _, l := range tenants {
		if l.ID == id {
			return l.Name
		}
	}

	return ""
}
