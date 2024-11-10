package graph_search

import "math/big"

// Bitset is a data structure that represents a bitset.
type Bitset struct {
	*big.Int
}

// NewBigInt returns a new instance of Bitset with an initialized big.Int.
func NewBigInt() Bitset {
	return Bitset{
		Int: new(big.Int),
	}
}

// Exists checks whether the bit at the specified index is set to 1.
func (b Bitset) Exists(i int32) bool {
	return b.Int.Bit(int(i)) == 1
}

// Set sets the value of the bit at the specified index to the given boolean value.
func (b Bitset) Set(i int32, value bool) {
	if value {
		b.Int.SetBit(b.Int, int(i), 1)
	} else {
		b.Int.SetBit(b.Int, int(i), 0)
	}
}

// Len returns the length of the bitset, which is the position of the highest bit set to 1.
func (b Bitset) Len() int {
	return b.Int.BitLen()
}
