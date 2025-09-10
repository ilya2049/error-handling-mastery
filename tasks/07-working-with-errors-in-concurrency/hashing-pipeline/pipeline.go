package pipeline

import (
	"context"
	"errors"
	"io"
	"strconv"
	"sync"
)

const maxBatchSize = 4

var (
	errEmptyInput       = errors.New("no data in input")
	errNilChannel       = errors.New("nil channel in input")
	errEmptyTx          = errors.New("input must not contain empty transactions")
	errInvalidBatchSize = errors.New("invalid batch size")
)

type Transaction struct {
	ID int64
}

func (t Transaction) Hash() Hash {
	return newHash([]byte(strconv.FormatInt(t.ID, 10)))
}

type Block struct {
	Hash Hash
}

// Pipeline представляет собой последовательное (но конкуретное внутри) выполнение следующих функций:
//
//	batch -> hashTransactions -> sink -> mergeErrors
//
// Pipeline гарантирует своё завершение только после завершения всех входящих в него функций.
// При этом пайплайн завершается:
//   - при отмене входящего контекста (как следствие того, что от контекста завершатся функции пайплайна);
//   - при получении хотя бы одной ошибки из финального канала ошибок (в таком случае пайплайн требует от функций
//     досрочного завершения и ожидает их);
//   - при безошибочной обработке всех транзакций.
func Pipeline(ctx context.Context, batchSize int, out io.Writer, txs ...Transaction) error {
	drainErrChan := func(errs <-chan error) {
		for range errs {
		}
	}

	var errChans []<-chan error

	batchChan, batchErrChan, err := batch(ctx, batchSize, txs...)
	if err != nil {
		return err
	}

	errChans = append(errChans, batchErrChan)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	blockChan, blockErrChan, err := hashTransactions(ctx, batchChan)
	if err != nil {
		cancel()

		outErrChan, _ := mergeErrors(errChans...)
		drainErrChan(outErrChan)

		return err
	}

	errChans = append(errChans, blockErrChan)

	sinkErrChan, err := sink(ctx, blockChan, out)
	if err != nil {
		cancel()

		outErrChan, _ := mergeErrors(errChans...)
		drainErrChan(outErrChan)

		return err
	}

	errChans = append(errChans, sinkErrChan)

	outErrChan, err := mergeErrors(errChans...)
	if err != nil {
		cancel()

		drainErrChan(outErrChan)

		return err
	}

	for err := range outErrChan {
		cancel()

		drainErrChan(outErrChan)

		return err
	}

	return nil
}

// batch преобразует входящий слайс транзакций в группы размером batchSize и отправляет в выходной канал:
//   - возвращает ошибку errInvalidBatchSize, если batchSize меньше 1 или больше maxBatchSize;
//   - возвращает ошибку errEmptyInput, если слайс транзакций пуст;
//   - если у очередной транзакции нулевой ID, то функция пишет в выходной канал ошибку errEmptyTx
//     и завершает своё выполнение.
func batch(ctx context.Context, batchSize int, txs ...Transaction) (<-chan []Transaction, <-chan error, error) {
	if batchSize < 1 || batchSize > maxBatchSize {
		return nil, nil, errInvalidBatchSize
	}

	if len(txs) == 0 {
		return nil, nil, errEmptyInput
	}

	txsChan := make(chan []Transaction)
	errChan := make(chan error, 1)

	go func() {
		batch := make([]Transaction, 0, batchSize)

		defer close(txsChan)
		defer close(errChan)

		for _, tx := range txs {
			if tx.ID == 0 {
				errChan <- errEmptyTx

				return
			}

			batch = append(batch, tx)
			if len(batch) == cap(batch) {
				select {
				case <-ctx.Done():
					return
				case txsChan <- batch:
					batch = make([]Transaction, 0, batchSize)
				}
			}
		}

		if len(batch) > 0 {
			select {
			case <-ctx.Done():
				return
			case txsChan <- batch:
			}
		}
	}()

	return txsChan, errChan, nil
}

// hashTransactions берёт группы транзакций из входного канала, и считает хеш от группы с помощью CalculateHash:
//   - возвращает ошибку errNilChannel, если на вход получила nil-канал;
//   - если при просчёте хеша возникает ошибка, то функция пишет её в выходной канал и завершает своё выполнение.
func hashTransactions(ctx context.Context, batchc <-chan []Transaction) (<-chan Block, <-chan error, error) {
	if batchc == nil {
		return nil, nil, errNilChannel
	}

	errChan := make(chan error, 2)
	blockc := make(chan Block)

	txsToHashable := func(txs []Transaction) []Hashable {
		hashable := make([]Hashable, 0, len(txs))

		for _, tx := range txs {
			hashable = append(hashable, tx)
		}

		return hashable
	}

	go func() {
		defer close(blockc)
		defer close(errChan)

		for batch := range batchc {
			hash, err := CalculateHash(txsToHashable(batch))
			if err != nil {
				errChan <- err

				return
			}

			select {
			case <-ctx.Done():
				return
			case blockc <- Block{Hash: hash}:
			}
		}
	}()

	return blockc, errChan, nil
}

// sink выводит в out хеши блоков из входного канала в строковой форме, разделяя их через '\n':
//   - возвращает ошибку errNilChannel, если на вход получила nil-канал;
//   - если при записи в out возникает ошибка, то функция пишет её в выходной канал и завершает своё выполнение.
func sink(ctx context.Context, blockc <-chan Block, out io.Writer) (<-chan error, error) {
	if blockc == nil {
		return nil, errNilChannel
	}

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		for block := range blockc {
			select {
			case <-ctx.Done():
				return
			default:
				outString := block.Hash.String() + "\n"
				if _, err := out.Write([]byte(outString)); err != nil {
					errChan <- err

					return
				}
			}
		}
	}()

	return errChan, nil
}

// mergeErrors сливает все ошибки из входящих каналов в выходной канал:
//   - возвращает ошибку errEmptyInput, если слайс каналов пуст;
//   - возвращает ошибку errNilChannel, если хотя бы один из каналов в слайсе нулевой;
//   - выходной канал закрывается только после того, как будут вычитаны все входные каналы.
func mergeErrors(errcs ...<-chan error) (<-chan error, error) {
	if len(errcs) == 0 {
		return nil, errEmptyInput
	}

	for _, ch := range errcs {
		if ch == nil {
			return nil, errNilChannel
		}
	}

	var wg sync.WaitGroup

	output := make(chan error)

	for _, ch := range errcs {
		wg.Add(1)
		go func(ch <-chan error) {
			for err := range ch {
				output <- err
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output, nil
}
