package checks

// OctopusClientErrorHandler is a service used to respond to errors from the octopus client.
type OctopusClientErrorHandler interface {
	// HandleError either handles the error and returns a OctopusCheckResult, or bubbles an error up
	HandleError(id string, group string, err error) (OctopusCheckResult, error)
}
