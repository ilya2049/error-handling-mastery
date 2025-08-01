package errs

func Extract(err error) []error {
	errs := unwrap(err)

	if len(errs) > 0 {
		return errs
	}

	return nil
}

type multipleUnwrapper interface {
	Unwrap() []error
}

type unwrapper interface {
	Unwrap() error
}

func unwrap(err error) []error {
	if err == nil {
		return []error{}
	}

	mult, ok := err.(multipleUnwrapper)
	if ok {
		var errs []error
		for _, e := range mult.Unwrap() {
			errs = append(errs, unwrap(e)...)
		}

		return errs
	}

	one, ok := err.(unwrapper)
	if ok {
		return unwrap(one.Unwrap())
	}

	return []error{err}
}
