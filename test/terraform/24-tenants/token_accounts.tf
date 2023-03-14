resource "octopusdeploy_token_account" "account_a" {
  description                       = "Account A"
  name                              = "Account A"
  environments                      = null
  tenant_tags                       = []
  tenants                           = [octopusdeploy_tenant.tenant_team_a.id, octopusdeploy_tenant.tenant_team_b.id]
  tenanted_deployment_participation = "TenantedOrUntenanted"
  token                             = "secretgoeshere"
}

resource "octopusdeploy_token_account" "account_b" {
  description                       = "Account B"
  name                              = "Account B"
  environments                      = null
  tenant_tags                       = []
  tenants                           = [octopusdeploy_tenant.tenant_team_a.id, octopusdeploy_tenant.tenant_team_b.id]
  tenanted_deployment_participation = "TenantedOrUntenanted"
  token                             = "secretgoeshere"
}