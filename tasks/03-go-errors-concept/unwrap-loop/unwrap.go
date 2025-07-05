package errs

type Unwrapper interface {
	Unwrap() error
}

func Unwrap(err error) error {
	for {
		u, ok := err.(Unwrapper)
		if !ok {
			return err
		}

		err = u.Unwrap()
	}
}
