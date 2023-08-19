package bloom

import "math"

const (
	seed = 0xbc9f1d34
	m    = 0xc6a4a793
)

// Filter is an encoded set of []byte keys.
type Filter []byte

// MayContainKey returns whether the filter may contain given key. False positives
func (f Filter) MayContainKey(k []byte) bool {
	return f.mayContain(Hash(k))
}

// MayContain returns whether the filter may contain given key. False positives
// are possible, where it returns true for keys not in the original set.
func (f Filter) mayContain(h uint32) bool {
	// check if the filter is empty
	if len(f) < 2 {
		return false
	}
	// obtain the number of hash functions
	k := f[len(f)-1]
	// if k > 30, this is reserved for potentially new encodings for short Bloom filters.
	if k > 30 {
		// This is reserved for potentially new encodings for short Bloom filters.
		// Consider it a match.
		return true
	}
	// calculate the total number of bits in the filter.
	nBits := uint32(8 * (len(f) - 1))
	// change the hash value by right shift and left shift to generate different bit positions for subsequent iterations.
	delta := h>>17 | h<<15
	for j := uint8(0); j < k; j++ {
		// For each hash function, calculate the bit position bitPos
		bitPos := h % nBits
		// Check if the corresponding bit has been set.
		// If the bit has not been set, the key is definitely not in the set, and false is returned.
		if f[bitPos/8]&(1<<(bitPos%8)) == 0 {
			return false
		}
		h += delta
	}
	return true
}

// NewFilter returns a new Bloom filter that encodes a set of []byte keys with
// the given number of bits per key, approximately.
//
// A good bitsPerKey value is 10, which yields a filter with ~ 1% false
// positive rate.
func NewFilter(keys []uint32, bitsPerKey int) Filter {
	return Filter(appendFilter(nil, keys, bitsPerKey))
}

// BloomBitsPerKey returns the bits per key required by bloomfilter based on
// the false positive rate.
func BloomBitsPerKey(numEntries int, fp float64) int {
	size := -1 * float64(numEntries) * math.Log(fp) / math.Pow(float64(0.69314718056), 2)
	locs := math.Ceil(float64(0.69314718056) * size / float64(numEntries))
	return int(locs)
}

func appendFilter(buf []byte, keys []uint32, bitsPerKey int) []byte {
	if bitsPerKey < 0 {
		bitsPerKey = 0
	}
	// 0.69 is approximately ln(2).
	k := uint32(float64(bitsPerKey) * 0.69)
	if k < 1 {
		k = 1
	}
	if k > 30 {
		k = 30
	}

	nBits := len(keys) * bitsPerKey
	// For small len(keys), we can see a very high false positive rate. Fix it
	// by enforcing a minimum bloom filter length.
	if nBits < 64 {
		nBits = 64
	}
	nBytes := (nBits + 7) / 8
	nBits = nBytes * 8
	buf, filter := extend(buf, nBytes+1)

	for _, h := range keys {
		delta := h>>17 | h<<15
		for j := uint32(0); j < k; j++ {
			bitPos := h % uint32(nBits)
			filter[bitPos/8] |= 1 << (bitPos % 8)
			h += delta
		}
	}
	filter[nBytes] = uint8(k)

	return buf
}

// extend appends n zero bytes to b. It returns the overall slice (of length
// n+len(originalB)) and the slice of n trailing zeroes.
func extend(b []byte, n int) (overall, trailer []byte) {
	want := n + len(b)
	if want <= cap(b) {
		overall = b[:want]
		trailer = overall[len(b):]
		for i := range trailer {
			trailer[i] = 0
		}
	} else {
		// Grow the capacity exponentially, with a 1KiB minimum.
		c := 1024
		for c < want {
			c += c / 4
		}
		overall = make([]byte, want, c)
		trailer = overall[len(b):]
		copy(overall, b)
	}
	return overall, trailer
}

// Hash implements a hashing algorithm similar to the Murmur hash.
func Hash(b []byte) uint32 {
	// The original algorithm uses a seed of 0x9747b28c.
	h := uint32(seed) ^ uint32(len(b))*m
	// Pick up four bytes at a time.
	for ; len(b) >= 4; b = b[4:] {
		// The original algorithm uses the following commented out code to load
		h += uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
		h *= m
		h ^= h >> 16
	}
	// Pick up remaining bytes.
	switch len(b) {
	case 3:
		h += uint32(b[2]) << 16
		fallthrough
	case 2:
		h += uint32(b[1]) << 8
		fallthrough
	case 1:
		h += uint32(b[0])
		h *= m
		h ^= h >> 24
	}
	return h
}
