package bloom

import (
	"github.com/spaolacci/murmur3"
	"math"
)

// Filter represents a structure for the filter itself.
type Filter struct {
	bitSet    []bool // Bit array to hold the state of the data
	size      uint32 // Size of the bit array
	numHashes uint8  // Number of hash functions to use
}

// NewBloomFilter initializes a new Bloom filter based on the expected number of items and desired false positive rate.
func NewBloomFilter(expectedItems uint32, fpRate float64) *Filter {
	// Calculate the size of bit array using the expected number of items and desired false positive rate
	size := uint32(-float64(expectedItems) * math.Log(fpRate) / (math.Ln2 * math.Ln2))
	// Calculate the optimal number of hash functions based on the size of bit array and expected number of items
	numHashes := uint8(float64(size) / float64(expectedItems) * math.Ln2)

	return &Filter{
		bitSet:    make([]bool, size),
		size:      size,
		numHashes: numHashes,
	}
}

// Add inserts an item into the Bloom filter.
func (f *Filter) Add(item []byte) {
	hashes := f.hash(item)
	// For each hash value, find the position and set the bit to true
	for i := uint8(0); i < f.numHashes; i++ {
		position := hashes[i] % f.size
		f.bitSet[position] = true
	}
}

// MayContainItem checks if an item is possibly in the set.
// If it returns false, the item is definitely not in the set.
// If it returns true, the item might be in the set, but it can also be a false positive.
func (f *Filter) MayContainItem(item []byte) bool {
	hashes := f.hash(item)
	for i := uint8(0); i < f.numHashes; i++ {
		position := hashes[i] % f.size
		if !f.bitSet[position] {
			return false
		}
	}
	return true
}

// hash produces multiple hash values for an item.
// It leverages two hash values from murmur3 and generates as many as needed through a linear combination.
func (f *Filter) hash(item []byte) []uint32 {
	h1, h2 := murmur3.Sum128(item) // Get two 64-bit hash values
	var result []uint32

	// Use the two hash values to generate the required number of hash functions.
	for i := uint8(0); i < f.numHashes; i++ {
		h := h1 + uint64(i)*h2
		result = append(result, uint32(h))
	}
	return result
}
