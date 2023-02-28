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

| Check ID                             | Description                                                |
|--------------------------------------|------------------------------------------------------------|
| OctoRecEnvironmentCount              | Counts the number of environments in the space.            |
 | OctoRecDefaultProjectGroupChildCount | Counts the number of projects in the default project group |
