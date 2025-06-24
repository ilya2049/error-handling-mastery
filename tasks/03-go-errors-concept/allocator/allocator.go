package allocator

const (
	Admin          = 777
	MinMemoryBlock = 1024
)

type NotPermittedError struct{}

func (*NotPermittedError) Error() string {
	return "operation not permitted"
}

type ArgOutOfDomainError struct{}

func (*ArgOutOfDomainError) Error() string {
	return "numerical argument out of domain of func"
}

func Allocate(userID, size int) ([]byte, error) {
	if userID != Admin {
		return nil, &NotPermittedError{}
	}

	if size < MinMemoryBlock {
		return nil, &ArgOutOfDomainError{}
	}

	return make([]byte, size), nil
}
