package checks

type OctopusCheck interface {
	Execute() (OctopusCheckResult, error)
}
