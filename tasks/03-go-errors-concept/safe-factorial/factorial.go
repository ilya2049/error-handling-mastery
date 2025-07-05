package factorial

import (
	"errors"
	"fmt"
)

const maxDepth = 256

// Реализуй нас.
var (
	ErrNegativeN = errors.New("negative n")
	ErrTooDeep   = fmt.Errorf("too deep, max depth is %d", maxDepth)
)

// Calculate рекурсивно считает факториал входного числа n.
// Если число меньше нуля, то возвращается ошибка ErrNegativeN.
// Если для вычисления факториала потребуется больше maxDepth фреймов, то Calculate вернёт ErrTooDeep.
func Calculate(n int) (int, error) {
	if n < 0 {
		return 0, ErrNegativeN
	}

	if n > maxDepth {
		return 0, ErrTooDeep
	}

	factorial := 1

	for i := 1; i <= n; i++ {
		factorial *= i
	}

	return factorial, nil
}
