package executor

import (
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"testing"
)

type alwaysFailCheck struct {
}

func (o alwaysFailCheck) Execute() (checks.OctopusCheckResult, error) {
	return checks.NewOctopusCheckResultImpl("This check always fails", "OctopusRecommendationAlwaysFail", "", checks.Error, ""), nil
}

type alwaysPassCheck struct {
}

func (o alwaysPassCheck) Execute() (checks.OctopusCheckResult, error) {
	return checks.NewOctopusCheckResultImpl("This check passed ok", "OctopusRecommendationAlwaysPass", "", checks.Ok, ""), nil
}

func TestNoChecks(t *testing.T) {
	results, err := OctopusCheckExecutor{}.ExecuteChecks(nil)

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if len(results) != 0 {
		t.Fatal("Should not have returned any results")
	}
}

func TestFailChecks(t *testing.T) {
	results, err := OctopusCheckExecutor{}.ExecuteChecks([]checks.OctopusCheck{alwaysFailCheck{}})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if len(results) != 1 {
		t.Fatal("Should have returned 1 results")
	}
}

func TestFailAndPassChecks(t *testing.T) {
	results, err := OctopusCheckExecutor{}.ExecuteChecks([]checks.OctopusCheck{alwaysFailCheck{}, alwaysPassCheck{}})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if len(results) != 2 {
		t.Fatal("Should have returned 2 results")
	}
}
