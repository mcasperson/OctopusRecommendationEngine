package executor

import "github.com/mcasperson/OctopusRecommendationEngine/internal/checks"

type OctopusCheckExecutor struct {
}

func (o OctopusCheckExecutor) ExecuteChecks(checkCollection []checks.OctopusCheck) ([]checks.OctopusCheckResult, error) {
	if checkCollection == nil || len(checkCollection) == 0 {
		return []checks.OctopusCheckResult{}, nil
	}

	checkResults := []checks.OctopusCheckResult{}

	for _, c := range checkCollection {
		result, err := c.Execute()

		if err != nil {
			return nil, err
		}

		checkResults = append(checkResults, result)
	}

	return checkResults, nil
}
