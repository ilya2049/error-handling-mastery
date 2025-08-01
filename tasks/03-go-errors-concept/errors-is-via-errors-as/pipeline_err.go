package pipe

import (
	"errors"
	"fmt"
)

type PipelineError struct {
	User        string
	Name        string
	FailedSteps []string
}

func (p *PipelineError) Error() string {
	return fmt.Sprintf("pipeline %q error", p.Name)
}

func IsPipelineError(err error, user, pipelineName string) bool {
	var pipelineError *PipelineError

	ok := errors.As(err, &pipelineError)

	if !ok {
		return false
	}

	return pipelineError.User == user && pipelineError.Name == pipelineName
}
