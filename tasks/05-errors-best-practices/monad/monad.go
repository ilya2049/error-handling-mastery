package monad

import "errors"

var ErrNoMonadValue = errors.New("no monad value")

// M представляет собой монаду.
type M struct {
	err error
	v   any
}

// Bind применяет функцию f к значению M, возвращая новую монаду.
// Если M невалидна, то Bind эффекта не имеет.
func (m M) Bind(f func(v any) M) M {
	if m.err != nil {
		return m
	}

	return f(m.v)
}

// Unpack возвращает значение и ошибку, хранимые в монаде.
// При отсутствии и ошибки и значения метод возвращает ErrNoMonadValue.
func (m M) Unpack() (any, error) {
	if (m == M{}) {
		return nil, ErrNoMonadValue
	}

	return m.v, m.err
}

// Unit конструирует M на основе значения v.
func Unit(v any) M {
	return M{
		v: v,
	}
}

// Err конструирует "невалидную" монаду M.
func Err(err error) M {
	return M{
		err: err,
	}
}
