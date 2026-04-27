package env

import (
	"math/rand"
	"sort"

	"github.com/subtlepseudonym/envcmp/internal/envfile"
)

// SampleOptions controls how Sample selects entries.
type SampleOptions struct {
	// N is the maximum number of entries to return.
	// If N <= 0 or N >= len(entries), all entries are returned.
	N int

	// Seed is used to initialize the random source for reproducibility.
	// If Seed == 0, selection is non-deterministic.
	Seed int64

	// Sorted causes the output slice to be sorted by key after sampling.
	Sorted bool
}

// Sample returns a random subset of up to N entries from the provided map.
// The original map is not modified.
func Sample(entries []envfile.Entry, opts SampleOptions) []envfile.Entry {
	if len(entries) == 0 {
		return []envfile.Entry{}
	}

	// Copy to avoid mutating caller's slice.
	pool := make([]envfile.Entry, len(entries))
	copy(pool, entries)

	if opts.N <= 0 || opts.N >= len(pool) {
		if opts.Sorted {
			sort.Slice(pool, func(i, j int) bool {
				return pool[i].Key < pool[j].Key
			})
		}
		return pool
	}

	var r *rand.Rand
	if opts.Seed != 0 {
		//nolint:gosec // reproducible sampling, not cryptographic use
		r = rand.New(rand.NewSource(opts.Seed))
	} else {
		//nolint:gosec
		r = rand.New(rand.NewSource(rand.Int63()))
	}

	// Fisher-Yates partial shuffle to pick N items.
	for i := 0; i < opts.N; i++ {
		j := i + r.Intn(len(pool)-i)
		pool[i], pool[j] = pool[j], pool[i]
	}

	result := pool[:opts.N]

	if opts.Sorted {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Key < result[j].Key
		})
	}

	return result
}
