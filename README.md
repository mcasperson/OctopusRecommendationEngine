# Octopus Linter

This CLI tool scans an Octopus instance to find potential issues in the configuration and suggests solutions.

## Usage

Download the latest binary from the [releases](https://github.com/mcasperson/OctopusRecommendationEngine/releases/latest).

```
./octolint \
    -apiKey API-YOURAPIKEY \
    -url https://yourinstance.octopus.app \
    -space Spaces-1234
```

## Example output

This is an example of the tool output:

```
[OctoRecDefaultProjectGroupChildCount] The default project group contains 79 projects. You may want to organize these projects into additional project groups.
[OctoRecEmptyProject] The following projects have no runbooks and no deployment process: Azure Octopus test, CloudFormation, K8s Yaml Import, AA Training Demo, Test2, Cac Vars, Helm Demo, Helm, Package, K8s
[OctoRecUnusedVariables] The following variables are unused: App Runner/x, App Runner/thisisnotused, Vars Demo/workerpool, GAE Node.js/Email Address, CloudFormation - Lambda/APIKey, ReleaseDiffTest/Variable2, ReleaseDiffTest/Variable2, ReleaseDiffTest/Variable1, Release Diff/packages[package].files[file].diff, Release Diff/scoped variable, Release Diff/scoped variable, Release Diff/unscoped variable, Var test/Config[Databases:Name].Value, Var test/Config[General:DEBUG].Value, K8s Command Example/MyTags[Three].Name, Rolling Deployment/b, Rolling Deployment/aws, Rolling Deployment/azure, Rolling Deployment/workerpool, Rolling Deployment/gcp, Rolling Deployment/cert, CloudFormation - RDS/APIKey, Devops Tasks/APIKey, Devops Tasks/AWS, Terraform Test/DockerHub.Password, Terraform Test/New Value, Terraform Test/Scoped value, CloudFormation - Lambda Simple/APIKey, AWS Account/AWS
[OctoRecDuplicatedVariables] The following variables are duplicated between projects. Consider moving these into library variable sets: Cart/MONGO_DB_HOSTNAME == Orders/MONGO_DB_HOSTNAME, Cart/MONGO_DB_HOSTNAME == User/MONGO_DB_HOSTNAME, Orders/MONGO_DB_HOSTNAME == User/MONGO_DB_HOSTNAME
[OctoRecTooManySteps] The following projects have 20 or more steps: K8s Yaml Import 2
```

## Checks

| Check ID (* Not Implemented Yet)                | Description                                                                             |
|-------------------------------------------------|-----------------------------------------------------------------------------------------|
| OctoRecEnvironmentCount                         | Counts the number of environments in the space.                                         |
 | OctoRecDefaultProjectGroupChildCount            | Counts the number of projects in the default project group.                             |
 | OctoRecEmptyProject                             | Finds projects with no deployment process and no runbooks.                              |
 | OctoRecProjectSpecificEnvs *                    | Finds environments that are specific to a single project.                               |
| OctoRecUnusedVariables                          | Finds unused variables in a project.                                                    |
 | OctoRecDuplicatedVariables                      | Finds variables with duplicated values.                                                 |
 | OctoRecAdminDeployments *                       | Finds deployments initiated by someone with admin credentials.                          |
 | OctoRecPerpetualApiKeys *                       | Finds API keys that do not expire.                                                      |
 | OctoRecUnusedApiKeys *                          | Finds API keys that have not been used in 30 days.                                      |
 | OctoRecSinglePhaseLifecycle *                   | Finds lifecycles with a single phase.                                                   |
 | OctoRecProjectGroupsWithExclusiveEnvironments * | Finds project groups with projects that have no common environments.                    |
| OctoRecSharedCloudAccounts *                    | Finds accounts that are shared between project groups.                                  |
 | OctoRecTooManySteps                             | Finds projects with too many deployment steps.                                          |
| OctoRecDirectTenantReferences *                 | Finds projects that reference common groups of tenants directly rather than using tags. |
 | OctoRecUnhealthyTargets *                       | Finds targets that have not been healthy in the last 30 days.                           |
 | OctoRecSharedGitUsername *                      | Finds projects that share git credentials.                                              |
 | OctoRecHitTaskCap *                             | Counts how many times Octopus hit the task cap.                                         |