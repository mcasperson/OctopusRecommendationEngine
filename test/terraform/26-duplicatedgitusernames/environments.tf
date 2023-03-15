resource "octopusdeploy_environment" "development_environment" {
  allow_dynamic_infrastructure = true
  description                  = "A development environment"
  name                         = "Development"
  use_guided_failure           = false
}

resource "octopusdeploy_environment" "test_environment" {
  allow_dynamic_infrastructure = true
  description                  = "A test environment"
  name                         = "Test"
  use_guided_failure           = false
}

resource "octopusdeploy_environment" "production_environment" {
  allow_dynamic_infrastructure = true
  description                  = "A production environment"
  name                         = "Production"
  use_guided_failure           = false
}

resource "octopusdeploy_environment" "environment1" {
  allow_dynamic_infrastructure = true
  description                  = "Another environment"
  name                         = "Environment1"
  use_guided_failure           = false
}

resource "octopusdeploy_environment" "environment2" {
  allow_dynamic_infrastructure = true
  description                  = "Another environment"
  name                         = "Environment2"
  use_guided_failure           = false
}

resource "octopusdeploy_environment" "environment3" {
  allow_dynamic_infrastructure = true
  description                  = "Another environment"
  name                         = "Environment3"
  use_guided_failure           = false
}
