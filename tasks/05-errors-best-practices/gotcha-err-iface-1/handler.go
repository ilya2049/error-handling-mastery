package rest

func Handle() error {
	if err := usefulWork(); err != nil {
		return ErrInternalServerError
	}

	return nil
}

var usefulWork = func() error {
	return nil
}
