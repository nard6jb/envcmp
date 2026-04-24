package env

// IntersectOptions controls the behaviour of Intersect.
type IntersectOptions struct {
	// MaskSecrets replaces sensitive values with "***" in the output.
	MaskSecrets bool
}

// IntersectEntry holds the values for a key that appears in both maps.
type IntersectEntry struct {
	Key    string
	Left   string
	Right  string
	Masked bool
}

// Intersect returns the set of keys present in both left and right, along with
// their respective values. Keys that appear in only one map are omitted.
func Intersect(left, right map[string]string, opts IntersectOptions) []IntersectEntry {
	var entries []IntersectEntry

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			continue
		}

		entry := IntersectEntry{
			Key:   k,
			Left:  lv,
			Right: rv,
		}

		if opts.MaskSecrets {
			entry.Left, entry.Right, entry.Masked = maskEntry(k, lv, rv)
		}

		entries = append(entries, entry)
	}

	sortEntries(entries)
	return entries
}

// sortEntries sorts IntersectEntry slice by Key for deterministic output.
func sortEntries(entries []IntersectEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Key < entries[j-1].Key; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
