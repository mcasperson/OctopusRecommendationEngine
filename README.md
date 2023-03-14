# Octolint

![image](https://user-images.githubusercontent.com/160104/222631936-e1ec480e-abd5-4622-978d-08259844aa14.png)

This CLI tool scans an Octopus instance to find potential issues in the configuration and suggests solutions.

## Support

This tool is **not** supported by Octopus. Feel free to report
an [issue](https://github.com/mcasperson/OctopusRecommendationEngine/issues).

This tool is also in an alpha state. We expect to add more checks and tweak the code before things stabalize.

## Usage

Download the latest binary from
the [releases](https://github.com/mcasperson/OctopusRecommendationEngine/releases/latest).

```
./octolint \
    -apiKey API-YOURAPIKEY \
    -url https://yourinstance.octopus.app \
    -space Spaces-1234
```

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

| Check ID (* Not Implemented Yet)                 | Description                                                                             |
|--------------------------------------------------|-----------------------------------------------------------------------------------------|
| OctoLintEnvironmentCount                         | Counts the number of environments in the space.                                         |
| OctoLintDefaultProjectGroupChildCount            | Counts the number of projects in the default project group.                             |
| OctoLintEmptyProject                             | Finds projects with no deployment process and no runbooks.                              |
| OctoLintProjectSpecificEnvs                      | Finds environments that are specific to a single project.                               |
| OctoLintUnusedVariables                          | Finds unused variables in a project.                                                    |
| OctoLintDuplicatedVariables                      | Finds variables with duplicated values.                                                 |
| OctoLintDeploymentQueuedByAdmin                  | Finds deployments initiated by someone with admin credentials.                          |
| OctoLintPerpetualApiKeys                         | Finds API keys that do not expire.                                                      |
| OctoLintProjectGroupsWithExclusiveEnvironments * | Finds project groups with projects that have no common environments.                    |
| OctoLintSharedCloudAccounts *                    | Finds accounts that are shared between project groups.                                  |
| OctoLintTooManySteps                             | Finds projects with too many deployment steps.                                          |
| OctoLintDirectTenantReferences                   | Finds projects that reference common groups of tenants directly rather than using tags. |
| OctoLintUnhealthyTargets *                       | Finds targets that have not been healthy in the last 30 days.                           |
| OctoLintSharedGitUsername *                      | Finds projects that share git credentials.                                              |
| OctoLintDeploymentQueuedTime                     | Counts how many times deployment tasks were queued for more than a minute.              |
| OctoLintUnusedTargets                            | Lists targets that have not performed a deployment in 30 days.                          |
