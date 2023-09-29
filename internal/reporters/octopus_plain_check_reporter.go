package reporters

import (
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"strings"
)

// OctopusPlainCheckReporter prints the lint reports in plain text to std out.
type OctopusPlainCheckReporter struct {
	minSeverity int
}

func NewOctopusPlainCheckReporter(minSeverity int) OctopusPlainCheckReporter {
	return OctopusPlainCheckReporter{minSeverity: minSeverity}
}

func (o OctopusPlainCheckReporter) Generate(results []checks.OctopusCheckResult) (string, error) {
	if results == nil || len(results) == 0 {
		return "", nil
	}

	report := []string{}

	for _, r := range results {
		if r.Severity() >= o.minSeverity {
			report = append(report, "========================================================================================================================")
			report = append(report, r.Code())
			report = append(report, r.Description())
		}
	}

	if len(report) == 0 {
		return "No issues detected", nil
	}

	return strings.Join(report[:], "\n\n"), nil
}
