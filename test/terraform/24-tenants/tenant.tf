resource "octopusdeploy_tenant" "tenant_team_a" {
  name        = "Team A"
  description = "Test tenant"
  tenant_tags = []
  depends_on = []

  project_environment {
    environments = [octopusdeploy_environment.test_environment.id, octopusdeploy_environment.development_environment.id, octopusdeploy_environment.production_environment.id]
    project_id   = octopusdeploy_project.deploy_frontend_project.id
  }
}

resource "octopusdeploy_tenant" "tenant_team_b" {
  name        = "Team B"
  description = "Test tenant"
  tenant_tags = []
  depends_on = []

  project_environment {
    environments = [octopusdeploy_environment.test_environment.id, octopusdeploy_environment.development_environment.id, octopusdeploy_environment.production_environment.id]
    project_id   = octopusdeploy_project.deploy_frontend_project.id
  }
}
