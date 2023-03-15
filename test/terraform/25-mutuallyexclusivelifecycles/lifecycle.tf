resource "octopusdeploy_lifecycle" "lifecycle1" {
  description = "This is the default lifecycle."
  name        = "Lifecycle 1"

  phase {
    optional_deployment_targets = [octopusdeploy_environment.development_environment.id]
    name                        = "Development"
  }

  phase {
    optional_deployment_targets = [octopusdeploy_environment.test_environment.id]
    name                        = "Test"
  }

  phase {
    optional_deployment_targets = [octopusdeploy_environment.production_environment.id]
    name                        = "Production"
  }

  depends_on  = [octopusdeploy_environment.development_environment, octopusdeploy_environment.test_environment, octopusdeploy_environment.production_environment]
}

resource "octopusdeploy_lifecycle" "lifecycle2" {
  description = "This is the default lifecycle."
  name        = "Lifecycle 2"

  phase {
    optional_deployment_targets = [octopusdeploy_environment.environment1.id]
    name                        = "Development"
  }

  phase {
    optional_deployment_targets = [octopusdeploy_environment.environment2.id]
    name                        = "Test"
  }

  phase {
    optional_deployment_targets = [octopusdeploy_environment.environment3.id]
    name                        = "Production"
  }

  depends_on  = [octopusdeploy_environment.environment1, octopusdeploy_environment.environment2, octopusdeploy_environment.environment3]
}