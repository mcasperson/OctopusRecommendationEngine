package checks

// OctopusClientErrorHandler is a service used to respond to errors from the octopus client.
type OctopusClientErrorHandler interface {
	// HandleError either handles the error and returns a OctopusCheckResult, or bubbles an error up. This
	// is called when the main request is made in a check, and the inability to process that first request
	// means the check can not continue.
	HandleError(id string, group string, err error) (OctopusCheckResult, error)
	// ShouldContinue is used to determine if a secondary request that failed should be ignored. This is called
	// when looping through resources where the failure of some requests do not invalidate the whole check. For example,
	// the inability to get a deployment process for one project does not invalidate the checks against other projects.
	ShouldContinue(err error) bool
}
