package checks

const (
	Error      int = 20
	Warning        = 15
	Info           = 10
	Permission     = 5
	Ok             = 0
)

const (
	Organization string = "Organization"
	Security            = "Security"
	Performance         = "Performance"
	Optimization        = "Optimization"
)

// OctopusCheckResult describes the result of an OctopusCheck
type OctopusCheckResult interface {
	Description() string
	Code() string
	Link() string
	Severity() int
	Category() string
}

type OctopusCheckResultImpl struct {
	description string
	code        string
	link        string
	severity    int
	category    string
}

func NewOctopusCheckResultImpl(description string, code string, link string, severity int, category string) OctopusCheckResultImpl {
	return OctopusCheckResultImpl{
		description: description,
		code:        code,
		link:        link,
		severity:    severity,
		category:    category,
	}
}

func (o OctopusCheckResultImpl) Description() string {
	return o.description
}

func (o OctopusCheckResultImpl) Code() string {
	return o.code
}

func (o OctopusCheckResultImpl) Link() string {
	return o.link
}

func (o OctopusCheckResultImpl) Severity() int {
	return o.severity
}

func (o OctopusCheckResultImpl) Category() string {
	return o.category
}
