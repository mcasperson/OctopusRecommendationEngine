resource "octopusdeploy_token_account" "account_a" {
  description                       = "Account A"
  name                              = "Account A"
  environments                      = null
  tenant_tags                       = ["tag1/a"]
  tenants                           = []
  tenanted_deployment_participation = "TenantedOrUntenanted"
  token                             = "secretgoeshere"
  depends_on                        = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]
}

resource "octopusdeploy_token_account" "account_b" {
  description                       = "Account B"
  name                              = "Account B"
  environments                      = null
  tenant_tags                       = ["tag1/a"]
  tenants                           = []
  tenanted_deployment_participation = "TenantedOrUntenanted"
  token                             = "secretgoeshere"
  depends_on                        = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]
}