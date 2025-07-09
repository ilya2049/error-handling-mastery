package pipe

import "fmt"

type UserError struct {
	Operation string
	User      string
}

func (u *UserError) Error() string {
	return fmt.Sprintf("user %s cannot do op %s", u.User, u.Operation)
}

type PipelineError struct {
	User        string
	Name        string
	FailedSteps []string
}

func (p *PipelineError) Error() string {
	return fmt.Sprintf("pipeline %q error", p.Name)
}

// Добавь метод As для типа *PipelineError.
func (p *PipelineError) As(obj any) bool {
	if obj == nil {
		return false
	}

	userErrorPtrPtr, ok := obj.(**UserError)
	if !ok {
		return false
	}

	if *userErrorPtrPtr == nil {
		*userErrorPtrPtr = &UserError{}
	}

	(*userErrorPtrPtr).Operation = p.Name
	(*userErrorPtrPtr).User = p.User

	return true
}
