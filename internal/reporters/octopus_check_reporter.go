package reporters

import "github.com/mcasperson/OctopusRecommendationEngine/internal/checks"

// OctopusCheckReporter defines the contract used by reporters to print the result of lint checks.
type OctopusCheckReporter interface {
	Generate(results []checks.OctopusCheckResult) (string, error)
}
