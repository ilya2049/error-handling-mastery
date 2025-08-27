package index

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIndexFromFileName(t *testing.T) {
	// Заполните кейсы ниже так, чтобы тест проходил.
	cases := []struct {
		fileName      string
		expectedIndex int
		expectedError error
	}{
		{
			fileName:      "parsed_page",
			expectedIndex: 0,
			expectedError: ErrInvalidFilename,
		},
		{
			fileName:      "parsedpage",
			expectedIndex: 0,
			expectedError: ErrInvalidFilename,
		},
		{
			fileName:      "parsed_page_",
			expectedIndex: 0,
			expectedError: ErrInvalidFilename,
		},
		{
			fileName:      "parsed_page_100_suffix",
			expectedIndex: 0,
			expectedError: strconv.ErrSyntax,
		},
		{
			fileName:      "parsed_page_-1",
			expectedIndex: 0,
			expectedError: ErrIndexMustBePositive,
		},
		{
			fileName:      "parsed_page_0",
			expectedIndex: 0,
			expectedError: ErrIndexMustBePositive,
		},
		{
			fileName:      "parsed_page_1",
			expectedIndex: 1,
			expectedError: nil,
		},
		{
			fileName:      "parsed_page_15.5",
			expectedIndex: 0,
			expectedError: strconv.ErrSyntax,
		},
		{
			fileName:      "parsed_page_1000",
			expectedIndex: 1000,
			expectedError: nil,
		},
		{
			fileName:      fmt.Sprintf("parsed_page_%d", math.MaxInt32+1),
			expectedIndex: 0,
			expectedError: strconv.ErrRange,
		},
		{
			fileName:      "absolutely incorrect file name",
			expectedIndex: 0,
			expectedError: ErrInvalidFilename,
		},
	}

	for _, tt := range cases {
		t.Run(tt.fileName, func(t *testing.T) {
			index, err := GetIndexFromFileName(tt.fileName)
			require.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedIndex, index)

			_, ok := tt.expectedError.(*strconv.NumError)
			assert.False(t, ok, "do not use *strconv.NumError directly, look for more specific errors")
		})
	}
}
