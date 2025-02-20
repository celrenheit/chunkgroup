package chunkgroup

import (
	"slices"
	"testing"

	"github.com/test-go/testify/require"
)

func TestBasic(t *testing.T) {
	var res [][]int

	cg := New(10, 1, func(t []int) error {
		res = append(res, t)
		return nil
	})

	for i := range 105 {
		cg.Add(i)
	}

	err := cg.Flush()
	require.NoError(t, err)

	var res2 []int
	for i := range 105 {
		res2 = append(res2, i)
	}

	expected := slices.Collect(slices.Chunk(res2, 10))
	require.Equal(t, expected, res)
}
