package pipe

import "fmt"

type PipelineError struct {
	User        string
	Name        string
	FailedSteps []string
}

func (p *PipelineError) Error() string {
	return fmt.Sprintf("pipeline %q error", p.Name)
}

// Добавь метод Is для типа *PipelineError.
func (p *PipelineError) Is(err error) bool {
	other, ok := err.(*PipelineError)
	if !ok {
		return false
	}

	return p.Name == other.Name && p.User == other.User
}
