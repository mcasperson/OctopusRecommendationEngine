# Octopus Recommendation Engine

This CLI tool scans an Octopus instance to find potential issues in the configuration and suggests solutions.

## Usage

```
./octorec \
    -apiKey API-YOURAPIKEY \
    -url https://yourinstance.octopus.app \
    -space Spaces-1234
```

## Checks

| Check ID (* Not Implemented Yet)                | Description                                                                             |
|-------------------------------------------------|-----------------------------------------------------------------------------------------|
| OctoRecEnvironmentCount                         | Counts the number of environments in the space.                                         |
 | OctoRecDefaultProjectGroupChildCount            | Counts the number of projects in the default project group.                             |
 | OctoRecEmptyProject                             | Finds projects with no deployment process and no runbooks.                              |
 | OctoRecProjectSpecificEnvs *                    | Finds environments that are specific to a single project.                               |
| OctoRecUnusedVariables *                        | Finds unused variables in a project.                                                    |
 | OctoRecDuplicatedVariables *                    | Finds variables with duplicated values.                                                 |
 | OctoRecAdminDeployments *                       | Finds deployments initiated by someone with admin credentials.                          |
 | OctoRecPerpetualApiKeys *                       | Finds API keys that do not expire.                                                      |
 | OctoRecUnusedApiKeys *                          | Finds API keys that have not been used in 30 days.                                      |
 | OctoRecSinglePhaseLifecycle *                   | Finds lifecycles with a single phase.                                                   |
 | OctoRecProjectGroupsWithExclusiveEnvironments * | Finds project groups with projects that have no common environments.                    |
| OctoRecSharedCloudAccounts *                    | Finds accounts that are shared between project groups.                                  |
 | OctoRecTooManySteps *                           | Finds projects with too many deployment steps.                                          |
| OctoRecDirectTenantReferences *                 | Finds projects that reference common groups of tenants directly rather than using tags. |
 | OctoRecUnhealthyTargets *                       | Finds targets that have not been healthy in the last 30 days.                           |
 | OctoRecSharedGitUsername *                      | Finds projects that share git credentials.                                              |
 | OctoRecHitTaskCap *                             | Counts how many times Octopus hit the task cap.                                         |