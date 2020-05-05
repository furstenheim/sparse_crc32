package sparse_crc32

import (
	"github.com/stretchr/testify/assert"
	"hash/crc32"
	"strconv"
	"testing"
)

func TestToSparseFile (t * testing.T) {
	tcs := []struct {
		input SparseFile
		expected []byte
	} {
		{
			SparseFile{[]SparseByte{}, 4},
			[]byte{0,0 ,0, 0},
		},		{
			SparseFile{[]SparseByte{{1, 2}}, 4},
			[]byte{0,2 ,0, 0},
		},{
			SparseFile{[]SparseByte{{1, 2}}, 5},
			[]byte{0,2 , 0, 0, 0},
		},
	}
	for _, tc := range tcs {
		res := tc.input.toBytes()
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
	}
}

func TestXORMultiply (t *testing.T) {
	tcs := []struct{
		input1, input2 uint32
		expected       uint64
	} {
		{1 << 3 + 1 << 2 + 1<<1, 1<< 3 + 1 << 2 + 1, 1 << 6 + 1<<2 + 1 << 1},
		{1 << 5, 1<<7, 1 << 12},
	}

	for _, tc := range tcs {
		res := xorMultiply32(tc.input1, tc.input2)
		assert.Equal(t, tc.expected, res)
	}
}


func TestCRCSparseFile (t *testing.T) {
	tcs := []struct {
		input    SparseFile
		expected []byte
	}{
		{
			SparseFile{[]SparseByte{}, 4},
			[]byte{0, 0, 0, 0},
		}, {
			SparseFile{[]SparseByte{{1, 2}}, 4},
			[]byte{0, 2, 0, 0},
		}, {
			SparseFile{[]SparseByte{{1, 2}}, 5},
			[]byte{0, 2, 0, 0, 0},
		}, {
			SparseFile{[]SparseByte{{0, 1}, {4, 55}}, 8},
			[]byte{1, 0, 0, 0, 55, 0, 0, 0},
		},
	}
	for _, tc := range tcs {
		res := tc.input.toBytes()
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
		assert.Equal(t, crc32.ChecksumIEEE(tc.expected), IEEESparse(tc.input))
	}
}

func TestReminderSparse (t *testing.T) {
	tcs := []struct {
		input    SparseFile
	}{
		{
			SparseFile{[]SparseByte{}, 4},
		}, {
			SparseFile{[]SparseByte{{3, 1}}, 4},
		}, {
			SparseFile{[]SparseByte{{1, 2}}, 5},
		}, {
			SparseFile{[]SparseByte{{1, 2}, {2, 3}}, 10},
		},
	}
	for _, tc := range tcs {
		res := tc.input.toBytes()
		assert.Equal(t, reminderIEEE(res), reminderIEEESparse(tc.input))
	}
}

func TestMultiplicative (t *testing.T) {
	increase1 := reminderIEEE([]byte{1, 0, 0, 0, 0})
	increase2 := reminderIEEE([]byte{1, 0, 0, 0, 0})
	increase3 := reminderIEEE([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, increase3, reminderIEEE(uint64ToArray(xorMultiply32(increase1, increase2))))
}


func TestPadded (t *testing.T) {
	increase1 := reminderIEEE([]byte{0, 1, 0, 0, 0, 0})
	increase2 := reminderIEEE([]byte{0,0 , 1, 0, 0, 0, 0})
	increase3 := reminderIEEE([]byte{0, 0, 0, 1, 0, 0, 0, 0})
	assert.Equal(t, increase2, increase3)
	assert.Equal(t, increase2, increase1)
}

func TestUint64ToArray (t * testing.T) {
	x, _ := strconv.ParseUint("101000101010000000101000101000101000001010001000100010100010101", 2, 64)
	tcs := []struct{
		expected []byte
		input uint64
	} {
		{[]byte{0,0, 0, 0, 0,0, 0, 0}, 0},
		{[]byte{0,0, 0, 0, 0,0, 1, 4}, 1 << 8 + 1 << 2},
		{[]byte{81, 80, 20, 81, 65, 68, 69, 21}, x},
	}

	for _, tc := range tcs {
		res := uint64ToArray(tc.input)
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
	}
}
