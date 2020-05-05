package sparse_crc32

import (
	"encoding/binary"
	"math/bits"
)

const (
	// IEEE is by far and away the most common CRC-32 polynomial.
	// Used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, ...
	IEEE = 0xedb88320
)

type SparseFile struct {
	FileBytes []SparseByte
	Size      uint64
}
type SparseByte struct {
	Position uint64
	Value    byte
}

// Computes the CRC32 for a file that is mostly null bytes, using multiplicative properties of taking mod.
func IEEESparse (file SparseFile) uint32 {
	position2Index := map[uint64]int{}
	for i , v := range(file.FileBytes) {
		file.FileBytes[i].Value = bits.Reverse8(v.Value)
		position2Index[v.Position] = i
	}
	for i := 0; i < 4; i++ {
		index, ok := position2Index[uint64(i)]
		if !ok {
			file.FileBytes = append(file.FileBytes, SparseByte{Position: uint64(i), Value: 0xFF})
		} else {
			file.FileBytes[index].Value ^= 0xFF
		}
	}

	// Add padding
	file.Size += 4
	newReminder := bits.Reverse32(reminderIEEESparse(file))

	return newReminder ^ 0xFFFFFFFF
}

func reminderIEEESparse(file SparseFile) uint32 {
	size := file.Size
	result := uint32(0)
	for _, sparseByte := range file.FileBytes {
		power := size - 1 - sparseByte.Position
		valueCoeff := reminderIEEE([]byte{0, 0, 0, sparseByte.Value})
		carry := uint32(1)
		power2Value := reminderIEEE([]byte{0, 0, 1, 0})
		for power > 0 {
			if power%2 == 1 {
				carry = multiplyOnModule(carry, power2Value)
			}
			power = power >> 1
			power2Value = multiplyOnModule(power2Value, power2Value)
		}
		result ^= multiplyOnModule(carry, valueCoeff)
	}

	return result
}

func reminderIEEE(array []byte) uint32{
	poly := bits.Reverse32(uint32(IEEE))
	current := (uint32(array[0]) << 24 + uint32(array[1]) << 16 + uint32(array[2]) << 8 + uint32(array[3]))
	if current & (1 << 31) == 1 {
		current ^= poly
	}
	for i := 0; i < len(array) - 4; i++ {
		nextIndex := i + 4;
		var nextValue byte
		nextValue = array[nextIndex]

		for i := 7; i >= 0; i-- {
			nextBit := (nextValue >> i) & 1
			if current >> 31 & 1 == 1 {
				current = ((current << 1) ^ uint32(nextBit)) ^ poly
			} else {
				current = (current << 1) ^ uint32(nextBit)
			}
		}
	}
	return current
}


func multiplyOnModule (a, b uint32) uint32 {
	return reminderIEEE(uint64ToArray(xorMultiply32(a, b)))
}

func xorMultiply32 (a, b uint32) uint64 {
	return xorMultiply(uint64(a), uint64(b))
}
func xorMultiply (a, b uint64) uint64{
	if a<1{return 0};return a%2*b^xorMultiply(a/2,b*2)
}

func uint64ToArray (a uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(a))
	return b
}


func (f SparseFile) toBytes () []byte {
	slice := make([]byte, f.Size)
	for _, v := range f.FileBytes {
		slice[v.Position] = v.Value
	}
	return slice
}
