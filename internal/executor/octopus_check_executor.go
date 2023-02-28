package executor

import (
	"github.com/avast/retry-go/v4"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
)

type OctopusCheckExecutor struct {
}

func NewOctopusCheckExecutor() OctopusCheckExecutor {
	return OctopusCheckExecutor{}
}

// ExecuteChecks executes each check and collects the results.
func (o OctopusCheckExecutor) ExecuteChecks(checkCollection []checks.OctopusCheck, handleError func(checks.OctopusCheck, error) error) ([]checks.OctopusCheckResult, error) {
	if checkCollection == nil || len(checkCollection) == 0 {
		return []checks.OctopusCheckResult{}, nil
	}

	checkResults := []checks.OctopusCheckResult{}

	for _, c := range checkCollection {
		err := retry.Do(
			func() error {
				result, err := c.Execute()

				if err != nil {
					return err
				}

				checkResults = append(checkResults, result)

				return nil
			}, retry.Attempts(3))

		if err != nil {
			err := handleError(c, err)
			if err != nil {
				return nil, err
			}
		}
	}

	return checkResults, nil
}
