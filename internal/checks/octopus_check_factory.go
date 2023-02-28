package checks

type OctopusCheckFactory struct {
}

func (o OctopusCheckFactory) BuildAllChecks() ([]OctopusCheck, error) {
	return []OctopusCheck{}, nil
}
