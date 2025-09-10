package pipeline

import "errors"

var errNothingToHash = errors.New("nothing to hash")

// CalculateHash реализует хеширование входящих элементов по принципу дерева хешей (Merkle tree).
// Если входящий слайс пуст, то возвращает ошибку errNothingToHash.
func CalculateHash(hh []Hashable) (Hash, error) {
	if len(hh) == 0 {
		return nil, errNothingToHash
	}

	hashes := make([]Hash, len(hh))
	for i, h := range hh {
		hashes[i] = h.Hash()
	}

	return calculateHash(hashes), nil
}

func calculateHash(hh []Hash) Hash {
	var newHashesLength int
	if len(hh)%2 == 0 {
		newHashesLength = len(hh) / 2
	} else {
		newHashesLength = len(hh)/2 + 1
	}

	newHashes := make([]Hash, 0, newHashesLength)

	for i := 0; i < len(hh); i += 2 {
		if i < len(hh)-1 {
			newHashes = append(newHashes, newHash(append(hh[i], hh[i+1]...)))
		} else {
			newHashes = append(newHashes, newHash(append(hh[i], hh[i]...)))
		}
	}

	if len(newHashes) > 1 {
		return calculateHash(newHashes)
	}

	return newHashes[0]
}
