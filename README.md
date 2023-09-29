# Octolint

![image](https://user-images.githubusercontent.com/160104/222631936-e1ec480e-abd5-4622-978d-08259844aa14.png)

This CLI tool scans an Octopus instance to find potential issues in the configuration and suggests solutions.

## Support

Feel free to report
an [issue](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/issues).

This tool is not officially supported by Octopus. Please do not contact the Octopus support channels regarding octolint.

## Usage

Download the latest binary from
the [releases](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/releases).

```
./octolint \
    -apiKey API-YOURAPIKEY \
    -url https://yourinstance.octopus.app \
    -space Spaces-1234
```

Octolint is also distributed as a Docker image:

```
docker run -t --rm octopussamples/octolint \
    -url https://yourinstance.octopus.app \
    -apiKey API-YOURAPIKEY \
    -space Spaces-1
```

## Permissions

`octolint` only requires read access - it does not modify anything on the server.

To create a read only account, deploy the Terraform module under the [serviceaccount](serviceaccount) directory:

```bash
export TF_VAR_octopus_server=https://yourinstance.octopus.app
export TF_VAR_octopus_apikey=API-apikeygoeshere
export TF_VAR_octopus_space_id=Spaces-#
cd serviceaccount
terraform init
terraform apply
```

This creates a user role, team, and service account all called `Octolint`. You can then create an API key for the service account, and use that API key with `octolint`. 

## Example output

This is an example of the tool output:

```
[OctoLintDefaultProjectGroupChildCount] The default project group contains 79 projects. You may want to organize these projects into additional project groups.
[OctoLintEmptyProject] The following projects have no runbooks and no deployment process: Azure Octopus test, CloudFormation, K8s Yaml Import, AA Training Demo, Test2, Cac Vars, Helm Demo, Helm, Package, K8s
[OctoLintUnusedVariables] The following variables are unused: App Runner/x, App Runner/thisisnotused, Vars Demo/workerpool, GAE Node.js/Email Address, CloudFormation - Lambda/APIKey, ReleaseDiffTest/Variable2, ReleaseDiffTest/Variable2, ReleaseDiffTest/Variable1, Release Diff/packages[package].files[file].diff, Release Diff/scoped variable, Release Diff/scoped variable, Release Diff/unscoped variable, Var test/Config[Databases:Name].Value, Var test/Config[General:DEBUG].Value, K8s Command Example/MyTags[Three].Name, Rolling Deployment/b, Rolling Deployment/aws, Rolling Deployment/azure, Rolling Deployment/workerpool, Rolling Deployment/gcp, Rolling Deployment/cert, CloudFormation - RDS/APIKey, Devops Tasks/APIKey, Devops Tasks/AWS, Terraform Test/DockerHub.Password, Terraform Test/New Value, Terraform Test/Scoped value, CloudFormation - Lambda Simple/APIKey, AWS Account/AWS
[OctoLintDuplicatedVariables] The following variables are duplicated between projects. Consider moving these into library variable sets: Cart/MONGO_DB_HOSTNAME == Orders/MONGO_DB_HOSTNAME, Cart/MONGO_DB_HOSTNAME == User/MONGO_DB_HOSTNAME, Orders/MONGO_DB_HOSTNAME == User/MONGO_DB_HOSTNAME
[OctoLintTooManySteps] The following projects have 20 or more steps: K8s Yaml Import 2
```

## Checks

| Check ID (* Not Implemented Yet)                                                                                                                                             | Description                                                                             |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------|
| [OctoLintEnvironmentCount](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintEnvironmentCount)                                             | Counts the number of environments in the space.                                         |
| [OctoLintDefaultProjectGroupChildCount](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintDefaultProjectGroupChildCount)                   | Counts the number of projects in the default project group.                             |
| [OctoLintEmptyProject](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintEmptyProject)                                                     | Finds projects with no deployment process and no runbooks.                              |
| [OctoLintProjectSpecificEnvs](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintProjectSpecificEnvs)                                       | Finds environments that are specific to a single project.                               |
| [OctoLintUnusedVariables](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintUnusedVariables)                                               | Finds unused variables in a project.                                                    |
| [OctoLintDuplicatedVariables](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintDuplicatedVariables)                                       | Finds variables with duplicated values.                                                 |
| [OctoLintDeploymentQueuedByAdmin](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintDeploymentQueuedByAdmin)                               | Finds deployments initiated by someone with admin credentials.                          |
| [OctoLintPerpetualApiKeys](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintPerpetualApiKeys)                                             | Finds API keys that do not expire.                                                      |
| [OctoLintProjectGroupsWithExclusiveEnvironments](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintProjectGroupsWithExclusiveEnvironments) | Finds project groups with projects that have no common environments.                    |
| [OctoLintTooManySteps](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintTooManySteps)                                                     | Finds projects with too many deployment steps.                                          |
| [OctoLintDirectTenantReferences](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintDirectTenantReferences)                                 | Finds projects that reference common groups of tenants directly rather than using tags. |
| [OctoLintUnhealthyTargets](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintUnhealthyTargets)                                             | Finds targets that have not been healthy in the last 30 days.                           |
| [OctoLintSharedGitUsername](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintSharedGitUsername)                                           | Finds projects that share git credentials.                                              |
| [OctoLintDeploymentQueuedTime](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintDeploymentQueuedTime)                                     | Counts how many times deployment tasks were queued for more than a minute.              |
| [OctoLintUnusedTargets](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoLintUnusedTargets)                                                   | Lists targets that have not performed a deployment in 30 days.                          |
| [OctoRecLifecycleRetention](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/wiki/OctoRecLifecycleRetention)                                           | Lists lifecycles that retain resources forever.                                         |
