package handmadestack

import (
	"errors"
	"fmt"
)

var (
	ErrExecSQL         = errors.New("exec sql error")
	ErrInitTransaction = errors.New("init transaction error")
)

type Entity struct {
	ID string
}

// Используются тестами.
var (
	getEntity        = func() (Entity, error) { return Entity{ID: "some-id"}, nil }
	updateEntity     = func(e Entity) error { return nil }
	runInTransaction = func(f func() error) error { return f() }
)

// Перепиши меня так, чтобы логика сохранилась,
// но путь до каждой ошибки был очевиден.
func handler() (Entity, error) {
	var e Entity

	if err := runInTransaction(func() (opErr error) {
		e, opErr = getEntity()
		if opErr != nil {
			return fmt.Errorf("getEntity: %v", opErr)
		}

		if err := updateEntity(e); err != nil {
			return fmt.Errorf("updateEntity 1: %v", err)
		}

		return nil
	}); err != nil {
		return Entity{}, fmt.Errorf("runInTransaction 1: %v", err)
	}

	if err := runInTransaction(func() error {
		if err := updateEntity(e); err != nil {
			return fmt.Errorf("updateEntity 2: %v", err)
		}

		return nil
	}); err != nil {
		return Entity{}, fmt.Errorf("runInTransaction 2: %v", err)
	}

	if err := runInTransaction(func() (opErr error) {
		if err := updateEntity(e); err != nil {
			return fmt.Errorf("updateEntity 3: %v", err)
		}

		return nil
	}); err != nil {
		return Entity{}, fmt.Errorf("runInTransaction 3: %v", err)
	}

	return e, nil
}
